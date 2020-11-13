package cmd

import (
	speech "cloud.google.com/go/speech/apiv1"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"videosubs/pkg/gcp"
)

func ProcessRemoteFile(remotePath string) error {

	speechCtx := context.Background()
	speechClient, err := speech.NewClient(speechCtx)
	if err != nil {
		return fmt.Errorf("error creating gcp speech client: %w", err)
	}

	log.Print("Getting recognize response")
	rresponse, err := gcp.GetRecognizeResponse(speechClient, remotePath)
	if err != nil {
		return fmt.Errorf("error getting recognize response: %w", err)
	}


	_, filename := filepath.Split(remotePath)
	ext := filepath.Ext(filename)
	basefilename := strings.TrimSuffix(filename, ext)

	err = SaveRecognizeResponse(rresponse, fmt.Sprintf("%v.json", basefilename))
	if err != nil {
		return fmt.Errorf("error saving recognize response to local files: %w", err)
	}

	log.Printf("Processing Recognize Response")
	err = ProcessRecognizeResponse(rresponse, basefilename)
	if err != nil {
		return fmt.Errorf("error processing recognize response: %w", err)
	}

	return nil
}

