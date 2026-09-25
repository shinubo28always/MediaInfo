// Unrated Coder t.me/Unrated_Coder

package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
)

func ExtractMediaInfo(ctx context.Context, url string) (map[string]interface{}, error) {
	if _, err := exec.LookPath("mediainfo"); err == nil {
		cmd := exec.CommandContext(ctx, "mediainfo", "--Output=JSON", url)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return nil, errors.New("mediainfo error: " + stderr.String())
		}

		var result map[string]interface{}
		if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
			return nil, err
		}
		return result, nil
	}

	if _, err := exec.LookPath("ffprobe"); err == nil {
		cmd := exec.CommandContext(ctx, "ffprobe",
			"-v", "quiet",
			"-print_format", "json",
			"-show_format",
			"-show_streams",
			url,
		)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return nil, errors.New("ffprobe error: " + stderr.String())
		}

		var result map[string]interface{}
		if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, errors.New("neither 'mediainfo' nor 'ffprobe' was found on the system")
}
