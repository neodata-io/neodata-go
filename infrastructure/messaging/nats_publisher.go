package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

// EventPublisher handles publishing events to NATS JetStream
type Publisher struct {
	jetStream  jetstream.JetStream
	retries    int
	retryDelay time.Duration
}

// NewEventPublisher creates a new instance of EventPublisher with optional retries and retry delay
func NewPublisher(client *NATSClient, retries int, retryDelay time.Duration) *Publisher {

	return &Publisher{
		jetStream:  client.js,
		retries:    retries,
		retryDelay: retryDelay,
	}
}

// Publish publishes an event to a specific subject
func (p *Publisher) Publish(ctx context.Context, subject string, data []byte) (*jetstream.PubAck, error) {
	// A helper method accepting subject and data as parameters
	ack, err := p.jetStream.Publish(ctx, subject, data)
	if err != nil {
		log.Printf("Failed to publish event to subject %s: %v", subject, err)
		return nil, err
	}

	fmt.Printf("Published msg with sequence number %d on stream %q", ack.Sequence, ack.Stream)

	return ack, nil
}
