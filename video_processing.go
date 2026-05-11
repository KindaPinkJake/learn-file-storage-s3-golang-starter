package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type ffprobeOutput struct {
	Streams []struct {
		DisplayAspectRatio string `json:"display_aspect_ratio"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffprobe error: %w", err)
	}

	var output ffprobeOutput
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		return "", fmt.Errorf("couldn't parse ffprobe output: %w", err)
	}

	if len(output.Streams) == 0 {
		return "", fmt.Errorf("no streams found in video")
	}

	switch output.Streams[0].DisplayAspectRatio {
	case "16:9":
		return "16:9", nil
	case "9:16":
		return "9:16", nil
	default:
		return "other", nil
	}
}

func processVideoForFastStart(filePath string) (string, error) {
	outputPath := filePath + ".processing"
	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputPath)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %w", err)
	}
	return outputPath, nil
}
