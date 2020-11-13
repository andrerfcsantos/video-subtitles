package videosubs

import "time"

type Transcript struct {
	Text       string           `json:"text"`
	Punctuated string           `json:"punctuated,omitempty"`
	Words      []TranscriptWord `json:"words"`
}

func (t *Transcript) AddWord(w TranscriptWord) {
	t.Words = append(t.Words, w)
}

type TranscriptWord struct {
	Word      string        `json:"word"`
	TimeStart time.Duration `json:"time_start"`
	TimeEnd   time.Duration `json:"time_end"`
}
