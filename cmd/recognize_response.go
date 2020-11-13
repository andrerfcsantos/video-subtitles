package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"videosubs/pkg/videosubs"

	log "github.com/sirupsen/logrus"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func ProcessRecognizeResponse(r *speechpb.LongRunningRecognizeResponse, name string) error {
	transcript := TranscriptFromRecognizeResponse(r)
	fileName := fmt.Sprintf("%v_transcript.json", name)

	log.Printf("Saving Transcript to %v", fileName)
	err := SaveTranscript(transcript, fileName)
	if err != nil {
		return fmt.Errorf("error saving transcript: %w", err)
	}

	log.Printf("Processing Transcript")
	err = ProcessTranscript(transcript, name)
	if err != nil {
		return fmt.Errorf("error processing transcript: %w", err)
	}

	return nil
}

func TranscriptFromRecognizeResponse(r *speechpb.LongRunningRecognizeResponse) *videosubs.Transcript {
	var builder strings.Builder
	var res videosubs.Transcript

	for _, result := range r.Results {
		for _, alt := range result.Alternatives {
			if alt.Transcript != "" {
				builder.WriteString(alt.Transcript)
				for _, w := range alt.Words {
					res.AddWord(videosubs.TranscriptWord{
						Word:      w.Word,
						TimeStart: time.Duration(w.StartTime.Seconds*1e9 + int64(w.StartTime.Nanos)),
						TimeEnd:   time.Duration(w.EndTime.Seconds*1e9 + int64(w.EndTime.Nanos)),
					})
				}
			}
		}
	}

	res.Text = builder.String()

	return &res
}

func LoadRecognizeResponse(filename string) (*speechpb.LongRunningRecognizeResponse, error) {
	var response speechpb.LongRunningRecognizeResponse

	responseData, err := ioutil.ReadFile(filename)
	if err != nil {
		return &response, fmt.Errorf("error reading json file: %w", err)
	}

	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return &response, fmt.Errorf("error unmarshaling data from json file: %w", err)
	}

	return &response, nil
}

func SaveRecognizeResponse(r *speechpb.LongRunningRecognizeResponse, filename string) error {
	f, err := os.Create(filename)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("error creating file %v: %w", filename, err)
	}

	jsonRes, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("error marshaling recognize response: %w", err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, jsonRes, "", "  ")
	if err != nil {
		return fmt.Errorf("error creating identation for file: %w", err)
	}

	_, err = out.WriteTo(f)
	if err != nil {
		return fmt.Errorf("error writting recognize response: %w", err)
	}

	return nil
}
