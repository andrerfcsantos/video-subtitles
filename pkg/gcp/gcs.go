package gcp

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func UploadFile(client *storage.Client, bucket, localFile, remoteFile string) error {
	ctx := context.Background()

	// Open local file.
	f, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	// Upload an object with storage.Writer.
	wc := client.Bucket(bucket).Object(remoteFile).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}

func EnsureBucketExists(client *storage.Client, bucketName, projectId string) error {
	exists, err := BucketExists(client, bucketName, projectId)
	if err != nil {
		return err
	}

	if !exists {
		err = CreateBucket(client, bucketName, projectId)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateBucket(client *storage.Client, bucketName, projectId string) error {
	ctx := context.Background()
	bkt := client.Bucket(bucketName)

	if err := bkt.Create(ctx, projectId, nil); err != nil {
		return err
	}
	return nil
}

func BucketExists(client *storage.Client, bucketName, projectId string) (bool, error) {
	res := false
	ctx := context.Background()
	bucketIt := client.Buckets(ctx, projectId)
	bucketIt.Prefix = bucketName

	for {
		battrs, err := bucketIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, err
		}
		if battrs.Name == bucketName {
			res = true
			break
		}
	}

	return res, nil
}
