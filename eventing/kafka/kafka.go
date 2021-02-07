package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Publisher is an implementation of the Publisher interface for Kafka.
type Publisher struct {
	w *kafka.Writer
}

// NewPublisher returns a new instance of Publisher for the given address and topic.
func NewPublisher(addr, topic string) *Publisher {
	return &Publisher{
		w: &kafka.Writer{
			Addr:     kafka.TCP(addr),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// Publish pushes a new message to Kafka, with the given key and
// a value as message as a JSON string.
func (p *Publisher) Publish(ctx context.Context, key string, message interface{}) error {
	bytes, _ := json.Marshal(message)
	m := kafka.Message{
		Key:   []byte(key),
		Value: bytes,
	}

	err := p.w.WriteMessages(ctx, m)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}
