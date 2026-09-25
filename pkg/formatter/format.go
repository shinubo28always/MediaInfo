// Unrated Coder t.me/Unrated_Coder

package formatter

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func getString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	val, ok := m[key]
	if !ok || val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func getNestedMap(m map[string]interface{}, key string) map[string]interface{} {
	if m == nil {
		return nil
	}
	val, ok := m[key]
	if !ok || val == nil {
		return nil
	}
	if nm, ok := val.(map[string]interface{}); ok {
		return nm
	}
	return nil
}

func getNestedSlice(m map[string]interface{}, key string) []interface{} {
	if m == nil {
		return nil
	}
	val, ok := m[key]
	if !ok || val == nil {
		return nil
	}
	if s, ok := val.([]interface{}); ok {
		return s
	}
	return nil
}

func align(key, value string, width int) string {
	if value == "" || value == "N/A" {
		return ""
	}
	if len(key) >= width {
		return fmt.Sprintf("%s : %s\n", key, value)
	}
	spaces := strings.Repeat(" ", width-len(key))
	return fmt.Sprintf("%s%s : %s\n", key, spaces, value)
}

func FormatOutput(data map[string]interface{}) string {
	if _, ok := data["media"]; ok {
		return formatMediaInfo(data)
	}
	if _, ok := data["streams"]; ok {
		return formatFFprobe(data)
	}
	return "❌ Could not parse media information."
}

func formatMediaInfo(data map[string]interface{}) string {
	var text strings.Builder
	media := getNestedMap(data, "media")
	tracks := getNestedSlice(media, "track")

	for _, trackItem := range tracks {
		track, ok := trackItem.(map[string]interface{})
		if !ok {
			continue
		}
		ttype := getString(track, "@type")

		switch ttype {
		case "General":
			text.WriteString("🗒 **General**\n")
			text.WriteString(align("Unique ID", getString(track, "UniqueID"), 40))
			text.WriteString(align("Format", getString(track, "Format"), 40))
			text.WriteString(align("Format version", getString(track, "Format_Version"), 40))

			fileSize := getString(track, "FileSize_String")
			if fileSize == "" {
				fileSize = getString(track, "FileSize")
			}
			text.WriteString(align("File size", fileSize, 40))

			duration := getString(track, "Duration_String3")
			if duration == "" {
				duration = getString(track, "Duration")
			}
			text.WriteString(align("Duration", duration, 40))
			text.WriteString(align("Overall bit rate", getString(track, "OverallBitRate_String"), 40))
			text.WriteString(align("Frame rate", getString(track, "FrameRate"), 40))
			text.WriteString(align("Movie name", getString(track, "Title"), 40))
			text.WriteString(align("Encoded date", getString(track, "Encoded_Date"), 40))
			text.WriteString(align("Writing application", getString(track, "Encoded_Application"), 40))
			text.WriteString(align("Writing library", getString(track, "Encoded_Library"), 40))
			text.WriteString("\n")

		case "Video":
			text.WriteString("🎞 **Video**\n")
			text.WriteString(align("ID", getString(track, "ID"), 40))
			text.WriteString(align("Format", getString(track, "Format"), 40))
			text.WriteString(align("Format/Info", getString(track, "Format_Info"), 40))
			text.WriteString(align("Format profile", getString(track, "Format_Profile"), 40))
			text.WriteString(align("Codec ID", getString(track, "CodecID"), 40))
			text.WriteString(align("Duration", getString(track, "Duration_String3"), 40))
			text.WriteString(align("Width", getString(track, "Width"), 40))
			text.WriteString(align("Height", getString(track, "Height"), 40))
			text.WriteString(align("Display aspect ratio", getString(track, "DisplayAspectRatio_String"), 40))
			text.WriteString(align("Frame rate mode", getString(track, "FrameRate_Mode_String"), 40))

			fr := getString(track, "FrameRate")
			if fr != "" {
				fr = fr + " FPS"
			}
			text.WriteString(align("Frame rate", fr, 40))
			text.WriteString(align("Color space", getString(track, "ColorSpace"), 40))
			text.WriteString(align("Chroma subsampling", getString(track, "ChromaSubsampling"), 40))

			bd := getString(track, "BitDepth")
			if bd != "" {
				bd = bd + " bits"
			}
			text.WriteString(align("Bit depth", bd, 40))
			text.WriteString(align("Title", getString(track, "Title"), 40))
			text.WriteString(align("Writing library", getString(track, "Encoded_Library"), 40))
			text.WriteString(align("Default", getString(track, "Default"), 40))
			text.WriteString(align("Forced", getString(track, "Forced"), 40))
			text.WriteString("\n")

		case "Audio":
			idxStr := getString(track, "StreamOrder")
			idx, _ := strconv.Atoi(idxStr)
			text.WriteString(fmt.Sprintf("🔊 **Audio #%d**\n", idx+1))
			text.WriteString(align("ID", getString(track, "ID"), 40))
			text.WriteString(align("Format", getString(track, "Format"), 40))
			text.WriteString(align("Codec ID", getString(track, "CodecID"), 40))
			text.WriteString(align("Duration", getString(track, "Duration_String3"), 40))

			ch := getString(track, "Channels")
			if ch != "" {
				ch = ch + " channels"
			}
			text.WriteString(align("Channel(s)", ch, 40))
			text.WriteString(align("Channel layout", getString(track, "ChannelLayout"), 40))
			text.WriteString(align("Sampling rate", getString(track, "SamplingRate_String"), 40))
			text.WriteString(align("Language", getString(track, "Language_String"), 40))
			text.WriteString(align("Title", getString(track, "Title"), 40))
			text.WriteString(align("Default", getString(track, "Default"), 40))
			text.WriteString(align("Forced", getString(track, "Forced"), 40))
			text.WriteString("\n")

		case "Text":
			idxStr := getString(track, "StreamOrder")
			idx, _ := strconv.Atoi(idxStr)
			text.WriteString(fmt.Sprintf("🔠 **Subtitle #%d**\n", idx+1))
			text.WriteString(align("ID", getString(track, "ID"), 40))
			text.WriteString(align("Format", getString(track, "Format"), 40))
			text.WriteString(align("Codec ID", getString(track, "CodecID"), 40))
			text.WriteString(align("Language", getString(track, "Language_String"), 40))
			text.WriteString(align("Title", getString(track, "Title"), 40))
			text.WriteString(align("Default", getString(track, "Default"), 40))
			text.WriteString(align("Forced", getString(track, "Forced"), 40))
			text.WriteString("\n")
		}
	}

	return fmt.Sprintf("```\n%s\n```", strings.TrimSpace(text.String()))
}

func formatFFprobe(data map[string]interface{}) string {
	var text strings.Builder
	format := getNestedMap(data, "format")
	streams := getNestedSlice(data, "streams")

	text.WriteString("🗒 **General**\n")
	text.WriteString(align("Format", getString(format, "format_long_name"), 40))

	sizeStr := getString(format, "size")
	if sizeStr != "" {
		size, _ := strconv.ParseFloat(sizeStr, 64)
		text.WriteString(align("File size", fmt.Sprintf("%.2f MB", size/(1024*1024)), 40))
	}

	durStr := getString(format, "duration")
	if durStr != "" {
		dur, _ := strconv.ParseFloat(durStr, 64)
		text.WriteString(align("Duration", fmt.Sprintf("%.2f min", dur/60.0), 40))
	}

	brStr := getString(format, "bit_rate")
	if brStr != "" {
		br, _ := strconv.ParseFloat(brStr, 64)
		text.WriteString(align("Bit rate", fmt.Sprintf("%.1f kbps", br/1000.0), 40))
	}
	text.WriteString("\n")

	for i, streamItem := range streams {
		stream, ok := streamItem.(map[string]interface{})
		if !ok {
			continue
		}
		stype := getString(stream, "codec_type")
		tags := getNestedMap(stream, "tags")

		switch stype {
		case "video":
			text.WriteString("🎞 **Video**\n")
			text.WriteString(align("Codec", getString(stream, "codec_name"), 40))

			widthStr := getString(stream, "width")
			heightStr := getString(stream, "height")
			if widthStr != "" && heightStr != "" {
				text.WriteString(align("Resolution", fmt.Sprintf("%sx%s", widthStr, heightStr), 40))
			}

			text.WriteString(align("Frame rate", getString(stream, "avg_frame_rate"), 40))

			bdStr := getString(stream, "bits_per_raw_sample")
			if bdStr != "" && bdStr != "0" {
				text.WriteString(align("Bit depth", bdStr+" bits", 40))
			}
			text.WriteString("\n")

		case "audio":
			text.WriteString(fmt.Sprintf("🔊 **Audio #%d**\n", i))
			text.WriteString(align("Codec", getString(stream, "codec_name"), 40))
			text.WriteString(align("Channels", getString(stream, "channels"), 40))
			if tags != nil {
				text.WriteString(align("Language", getString(tags, "language"), 40))
				text.WriteString(align("Title", getString(tags, "title"), 40))
			}
			text.WriteString("\n")

		case "subtitle":
			text.WriteString(fmt.Sprintf("🔠 **Subtitle #%d**\n", i))
			text.WriteString(align("Codec", getString(stream, "codec_name"), 40))
			if tags != nil {
				text.WriteString(align("Language", getString(tags, "language"), 40))
				text.WriteString(align("Title", getString(tags, "title"), 40))
			}
			text.WriteString("\n")
		}
	}

	// Unrated Coder t.me/Unrated_Coder
	return fmt.Sprintf("```\n%s\n```", strings.TrimSpace(text.String()))
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
