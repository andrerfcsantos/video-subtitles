package gcp

import (
	"context"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

func DefaultCredentials() (*google.Credentials, error) {
	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}
