package gcp

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"

	"github.com/reecerussell/open-social/media"
)

// bucket is an implementation of media.Service for Google Cloud Platform's storage buckets.
type bucket struct {
	bucketName     string
	timeoutSeconds int
}

// New returns a new instance of media.Service for GCP storeage buckets.
func New(bucketName string) media.Service {
	return &bucket{
		bucketName:     bucketName,
		timeoutSeconds: 30,
	}
}

func (b *bucket) Upload(ctx context.Context, key string, data []byte) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to init client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(b.timeoutSeconds))
	defer cancel()

	w := client.Bucket(b.bucketName).Object(key).NewWriter(ctx)
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	return nil
}

func (b *bucket) Download(ctx context.Context, key string) ([]byte, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to init client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(b.timeoutSeconds))
	defer cancel()

	r, err := client.Bucket(b.bucketName).Object(key).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("reader: %v", err)
	}
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read: %v", err)
	}

	return data, nil
}
