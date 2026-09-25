// Unrated Coder t.me/Unrated_Coder

package formatter

import (
	"strings"
	"testing"
)

func TestFormatMediaInfo(t *testing.T) {
	data := map[string]interface{}{
		"media": map[string]interface{}{
			"track": []interface{}{
				map[string]interface{}{
					"@type":            "General",
					"Format":           "Matroska",
					"FileSize_String":  "50.0 MiB",
					"Duration_String3": "00:05:00.000",
				},
				map[string]interface{}{
					"@type":    "Video",
					"Format":   "AVC",
					"Width":    "1920",
					"Height":   "1080",
					"BitDepth": "8",
				},
			},
		},
	}

	res := FormatOutput(data)
	if !strings.Contains(res, "General") {
		t.Errorf("Expected General block, got %s", res)
	}
	if !strings.Contains(res, "Matroska") {
		t.Errorf("Expected format Matroska, got %s", res)
	}
	if !strings.Contains(res, "Video") {
		t.Errorf("Expected Video block, got %s", res)
	}
	if !strings.Contains(res, "1920x1080") && !strings.Contains(res, "Width") {
		t.Errorf("Expected Width details, got %s", res)
	}
}

func TestFormatFFprobe(t *testing.T) {
	data := map[string]interface{}{
		"format": map[string]interface{}{
			"format_long_name": "Matroska / WebM",
			"size":             "52428800",
			"duration":         "300.000000",
			"bit_rate":         "1398101",
		},
		"streams": []interface{}{
			map[string]interface{}{
				"codec_type":          "video",
				"codec_name":          "h264",
				"width":               1920,
				"height":              1080,
				"avg_frame_rate":      "30/1",
				"bits_per_raw_sample": "8",
			},
		},
	}

	res := FormatOutput(data)
	if !strings.Contains(res, "General") {
		t.Errorf("Expected General block, got %s", res)
	}
	if !strings.Contains(res, "Matroska / WebM") {
		t.Errorf("Expected format Matroska / WebM, got %s", res)
	}
	if !strings.Contains(res, "50.00 MB") {
		t.Errorf("Expected File size 50.00 MB, got %s", res)
	}
	if !strings.Contains(res, "5.00 min") {
		t.Errorf("Expected Duration 5.00 min, got %s", res)
	}
}
