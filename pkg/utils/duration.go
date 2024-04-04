package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetVideoDuration(videoData []byte) (time.Duration, error) {
	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "null", "-")
	cmd.Stdin = bytes.NewReader(videoData)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("error running ffmpeg command: %v", err)
	}

	durationString := extractDuration(string(output))
	if durationString == "" {
		return 0, errors.New("failed to extreact duration")
	}

	parts := strings.Split(durationString, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("unexpected duration format: %s", durationString)
	}

	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])
	secondsWithMilliseconds := strings.Split(parts[2], ".")
	seconds, _ := strconv.Atoi(secondsWithMilliseconds[0])

	// Convert milliseconds to nanoseconds
	milliseconds, _ := strconv.Atoi(secondsWithMilliseconds[1])
	nanoseconds := milliseconds * 1000000

	// Construct the duration
	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second + time.Duration(nanoseconds)

	return duration, nil
}

func extractDuration(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Duration:") {
			parts := strings.Split(line, ",")
			for _, part := range parts {
				if strings.Contains(part, "Duration:") {
					return strings.TrimSpace(strings.Split(part, ": ")[1])
				}
			}
		}
	}
	return ""
}
