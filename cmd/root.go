package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var DefaultBucketName = "video-subtitles-2"

var rootCmd = &cobra.Command{
	Use:   "videosubs",
	Short: "Video Subtitles is a subtitle generator for videos.",
	Long:  `Automatically generate subtitles for your video based on its audio.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return fmt.Errorf("this command takes at least 1 argument with the video file to process or " +
				"the JSON file with the transcript results")
		}

		path := args[0]

		var err error
		if strings.HasPrefix(path, "gs://") {
			err = ProcessRemoteFile(path)
		} else {
			err = ProcessLocalFile(path)
		}

		if err != nil {
			return fmt.Errorf("error creating subtitles for file '%v': %w", path, err)
		}

		return nil
	},
}

func Execute() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "15:04:05 MST",
		FullTimestamp:   true,
		ForceColors:     true,
	})

	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
