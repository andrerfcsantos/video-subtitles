package gcp

import (
	"context"
	"fmt"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
)

func GetRecognizeResponse(client *speech.Client, gsUri string) (*speechpb.LongRunningRecognizeResponse, error) {

	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			LanguageCode:                        "en-US",
			MaxAlternatives:                     1,
			EnableWordTimeOffsets:               true,
			Model:                               "video",
			UseEnhanced:                         true,
			EnableAutomaticPunctuation:          true,
			ProfanityFilter:                     false,
			EnableSeparateRecognitionPerChannel: false,
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gsUri},
		},
	}

	ctx := context.Background()

	op, err := client.LongRunningRecognize(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error performing long recognize request: %w", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for long recognize operation to finish: %w", err)
	}

	return resp, nil
}
