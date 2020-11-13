package cmd

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"videosubs/pkg/ffmpeg"
	"videosubs/pkg/gcp"
)

func ProcessLocalFile(path string) error {
	dir, filename := filepath.Split(path)
	ext := filepath.Ext(filename)

	if ext == ".json" {
		return ProcessLocalJsonFile(path)
	}

	basefilename := strings.TrimSuffix(filename, ext)

	audioFileName := fmt.Sprintf("%v.wav", basefilename)
	audioFilePath := fmt.Sprintf("%v%v", dir, audioFileName)

	log.Printf("Extracting audio from video file")
	err := ffmpeg.Convert(path, audioFilePath)
	if err != nil {
		return err
	}

	ctx := context.Background()
	stClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error creating gcp storage client: %w", err)
	}

	defaultCredentials, err := gcp.DefaultCredentials()
	if err != nil {
		return fmt.Errorf("error getting gcp default credentials: %w", err)
	}

	log.Printf("Making sure bucket exists before upload")
	err = gcp.EnsureBucketExists(stClient, DefaultBucketName, defaultCredentials.ProjectID)
	if err != nil {
		return fmt.Errorf("error getting ensuring bucket '%v' exists: %w", DefaultBucketName, err)
	}

	log.Printf("Uploading file")
	err = gcp.UploadFile(stClient, DefaultBucketName, audioFilePath, audioFileName)
	if err != nil {
		return fmt.Errorf("error uploading audio file '%v' : %w", audioFilePath, err)
	}

	err = ProcessRemoteFile(fmt.Sprintf("gs://%v/%v", DefaultBucketName, audioFileName))
	if err != nil {
		return err
	}

	return nil
}

func ProcessLocalJsonFile(path string) error {

	_, filename := filepath.Split(path)
	ext := filepath.Ext(filename)
	basefilename := strings.TrimSuffix(filename, ext)

	log.Printf("Loading Recognize Response from %v", path)
	rresponse, err := LoadRecognizeResponse(path)
	if err != nil {
		return fmt.Errorf("error saving file transcript: %w", err)
	}

	log.Printf("Processing Recognize Response")
	err = ProcessRecognizeResponse(rresponse, basefilename)
	if err != nil {
		return fmt.Errorf("error processing recognize response: %w", err)
	}

	return nil
}