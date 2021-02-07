package media

import "context"

// Service is a high level interface used to manage media content.
type Service interface {
	Upload(ctx context.Context, key string, data []byte) error
	Download(ctx context.Context, key string) ([]byte, error)
}
