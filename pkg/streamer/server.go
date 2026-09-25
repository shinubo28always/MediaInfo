// Unrated Coder t.me/Unrated_Coder

package streamer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gotd/td/tg"
)

type DocumentInfo struct {
	ID            int64
	AccessHash    int64
	FileReference []byte
	Size          int64
	MimeType      string
	FileName      string
}

type Streamer struct {
	api    *tg.Client
	store  sync.Map
	server *http.Server
}

func NewStreamer(api *tg.Client) *Streamer {
	return &Streamer{
		api: api,
	}
}

func (s *Streamer) SetAPI(api *tg.Client) {
	s.api = api
}

func (s *Streamer) RegisterDocument(token string, doc *tg.Document) {
	var fileName string = "file"
	for _, attr := range doc.Attributes {
		if fnAttr, ok := attr.(*tg.DocumentAttributeFilename); ok {
			fileName = fnAttr.FileName
			break
		}
	}

	info := DocumentInfo{
		ID:            doc.ID,
		AccessHash:    doc.AccessHash,
		FileReference: doc.FileReference,
		Size:          doc.Size,
		MimeType:      doc.MimeType,
		FileName:      fileName,
	}
	s.store.Store(token, info)
}

func (s *Streamer) Start(host string, port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/stream/", s.handleStream)

	addr := fmt.Sprintf("%s:%d", host, port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Printf("Local streamer started at http://%s", addr)
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Streamer server error: %v", err)
		}
	}()
}

func (s *Streamer) Stop() error {
	if s.server != nil {
		return s.server.Shutdown(context.Background())
	}
	return nil
}

func (s *Streamer) handleStream(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	token := parts[2]

	val, ok := s.store.Load(token)
	if !ok {
		http.Error(w, "File not found or link expired", http.StatusNotFound)
		return
	}
	info := val.(DocumentInfo)

	fileSize := info.Size
	start := int64(0)
	end := fileSize - 1

	rangeHeader := r.Header.Get("Range")
	hasRange := false

	if rangeHeader != "" {
		re := regexp.MustCompile(`bytes=(\d+)-(\d*)`)
		matches := re.FindStringSubmatch(rangeHeader)
		if len(matches) >= 2 {
			hasRange = true
			parsedStart, err := strconv.ParseInt(matches[1], 10, 64)
			if err == nil {
				start = parsedStart
			}
			if len(matches) > 2 && matches[2] != "" {
				parsedEnd, err := strconv.ParseInt(matches[2], 10, 64)
				if err == nil {
					end = parsedEnd
				}
			}
		}
	}

	if start < 0 || start >= fileSize || end < start {
		http.Error(w, "Requested range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	if end >= fileSize {
		end = fileSize - 1
	}

	limit := int(end - start + 1)

	mimeType := info.MimeType
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, info.FileName))

	var status int
	if hasRange {
		status = http.StatusPartialContent
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	} else {
		status = http.StatusOK
	}

	w.Header().Set("Content-Length", strconv.Itoa(limit))
	w.WriteHeader(status)

	if r.Method == "HEAD" {
		return
	}

	location := &tg.InputDocumentFileLocation{
		ID:            info.ID,
		AccessHash:    info.AccessHash,
		FileReference: info.FileReference,
	}

	chunkOffset := start
	bytesLeft := int64(limit)

	for bytesLeft > 0 {
		alignedOffset := (chunkOffset / 4096) * 4096
		endPos := chunkOffset + bytesLeft
		alignedEnd := ((endPos + 4095) / 4096) * 4096
		alignedLimit := int(alignedEnd - alignedOffset)

		if alignedLimit > 1048576 {
			alignedLimit = 1048576
		}

		req := &tg.UploadGetFileRequest{
			Location: location,
			Offset:   alignedOffset,
			Limit:    alignedLimit,
		}

		res, err := s.api.UploadGetFile(r.Context(), req)
		if err != nil {
			log.Printf("UploadGetFile error: %v", err)
			return
		}

		switch f := res.(type) {
		case *tg.UploadFile:
			bytesOffset := chunkOffset - alignedOffset
			if bytesOffset >= int64(len(f.Bytes)) {
				return
			}
			sliceEnd := bytesOffset + bytesLeft
			if sliceEnd > int64(len(f.Bytes)) {
				sliceEnd = int64(len(f.Bytes))
			}
			chunkBytes := f.Bytes[bytesOffset:sliceEnd]
			_, err = w.Write(chunkBytes)
			if err != nil {
				return
			}
			chunkOffset += int64(len(chunkBytes))
			bytesLeft -= int64(len(chunkBytes))
		default:
			log.Printf("Unexpected UploadFile class type or redirect received")
			return
		}
	}
}
