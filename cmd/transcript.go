package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"videosubs/pkg/videosubs"

	log "github.com/sirupsen/logrus"

	astisub "github.com/asticode/go-astisub"
)

func ProcessTranscript(t *videosubs.Transcript, name string) error {

	var currentChunk []videosubs.TranscriptWord
	var chunks [][]videosubs.TranscriptWord

	maxGap := time.Duration(1e9 + 5e8)

	n := len(t.Words)
	currentChunkChars := 0

	for i := 0; i < n-1; i++ {
		w := t.Words[i]
		nw := t.Words[i+1]
		gap := nw.TimeStart - w.TimeEnd

		if gap > maxGap || (strings.HasSuffix(w.Word, ".") && currentChunkChars > 35) {
			currentChunk = append(currentChunk, w)
			chunks = append(chunks, currentChunk)
			currentChunk = []videosubs.TranscriptWord{}
			currentChunkChars = 0
			continue
		}

		if currentChunkChars > 65 {
			chunks = append(chunks, currentChunk)
			currentChunk = []videosubs.TranscriptWord{}
			currentChunkChars = 0
		}

		currentChunk = append(currentChunk, w)
		currentChunkChars += 1 + len(w.Word)
	}

	currentChunk = append(currentChunk, t.Words[n-1])
	chunks = append(chunks, currentChunk)

	subtitles := astisub.NewSubtitles()
	for i, chunk := range chunks {
		firstWord := chunk[0]
		lastWord := chunk[len(chunk)-1]

		var words []string
		chunkChars := 0

		for _, word := range chunk {
			words = append(words, word.Word)
			chunkChars += len(word.Word)
		}

		var lines []astisub.Line
		if chunkChars > 40 && len(words) > 1 {
			half := chunkChars / 2
			c := 0
			for i := 0; i < len(words)-1; i++ {
				c += len(words[i])
				if c+len(words[i+1]) > half {
					lines = []astisub.Line{
						{
							Items: []astisub.LineItem{
								{
									InlineStyle: nil,
									Style:       nil,
									Text:        strings.Join(words[:i+1], " "),
								},
							},
							VoiceName: "",
						},
						{
							Items: []astisub.LineItem{{
								InlineStyle: nil,
								Style:       nil,
								Text:        strings.Join(words[i+1:], " "),
							},
							},
							VoiceName: "",
						},
					}
					break
				}
			}
		} else {
			lines = []astisub.Line{
				{
					Items: []astisub.LineItem{
						{
							InlineStyle: nil,
							Style:       nil,
							Text:        strings.Join(words, " "),
						},
					},
					VoiceName: "",
				},
			}
		}

		item := astisub.Item{
			Comments:    nil,
			Index:       i,
			EndAt:       lastWord.TimeEnd,
			InlineStyle: nil,
			Lines:       lines,
			Region:      nil,
			StartAt:     firstWord.TimeStart,
			Style:       nil,
		}

		subtitles.Items = append(subtitles.Items, &item)
	}

	fileName := fmt.Sprintf("%v.srt", name)
	log.Printf("Writting subtitle to %v", fileName)
	f, err := os.Create(fileName)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("could not open subtitle file: %w", err)
	}

	err = (*subtitles).WriteToSRT(f)
	if err != nil {
		return fmt.Errorf("could not write subtitle file: %w", err)
	}

	return nil
}

func SaveTranscript(t *videosubs.Transcript, filename string) error {
	f, err := os.Create(filename)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("error creating file to save transcript: %w", err)
	}

	jsonRes, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("error marshaling processed transcript data: %w", err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, jsonRes, "", "  ")
	if err != nil {
		return fmt.Errorf("error creating identation for file: %w", err)
	}

	_, err = out.WriteTo(f)
	if err != nil {
		return fmt.Errorf("error wriitting processed transcript data to json file: %w", err)
	}

	return nil
}
