package eventing

import "context"

// Publisher is a high-level interface used to publish messages to a sub/pub.
type Publisher interface {
	Publish(ctx context.Context, key string, message interface{}) error
}
