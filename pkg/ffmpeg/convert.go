package ffmpeg

import (
	"fmt"
	"os/exec"
)

func Convert(input, output string) error {

	ffmpegCmd := exec.Command("ffmpeg", "-y", "-i", input, "-ac", "1", output)

	err := ffmpegCmd.Run()
	if err != nil {
		return fmt.Errorf("error running ffmpeg command: %w", err)
	}

	return nil
}