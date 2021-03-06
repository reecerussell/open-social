package gcp

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/reecerussell/open-social/media"
)

var credentialJSON = "GOOGLE_CREDENTIAL_JSON"

// bucket is an implementation of media.Service for Google Cloud Platform's storage buckets.
type bucket struct {
	bucketName     string
	timeoutSeconds int
	client         *storage.Client
}

// New returns a new instance of media.Service for GCP storeage buckets.
func New(ctx context.Context, bucketName string) (media.Service, error) {
	var client *storage.Client
	var err error

	credJSON := os.Getenv(credentialJSON)
	if credJSON != "" {
		cred := option.WithCredentialsJSON([]byte(credJSON))
		client, err = storage.NewClient(ctx, cred)
	} else {
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to init client: %v", err)
	}

	return &bucket{
		bucketName:     bucketName,
		timeoutSeconds: 30,
		client:         client,
	}, nil
}

func (b *bucket) Upload(ctx context.Context, key string, data []byte) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(b.timeoutSeconds))
	defer cancel()

	w := b.client.Bucket(b.bucketName).Object(key).NewWriter(ctx)
	_, err := w.Write(data)
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(b.timeoutSeconds))
	defer cancel()

	r, err := b.client.Bucket(b.bucketName).Object(key).NewReader(ctx)
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
