package messaging

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// EventPublisher handles publishing events to NATS JetStream
type EventPublisher struct {
	jetStream  nats.JetStreamContext
	retries    int
	retryDelay time.Duration
}

// NewEventPublisher creates a new instance of EventPublisher with optional retries and retry delay
func NewEventPublisher(client *NATSClient, retries int, retryDelay time.Duration) *EventPublisher {
	return &EventPublisher{
		jetStream:  client.JetStream,
		retries:    retries,
		retryDelay: retryDelay,
	}
}

// Publish publishes an event to a specific subject
func (p *EventPublisher) Publish(ctx context.Context, subject string, data []byte) (*nats.PubAck, error) {
	ack, err := p.jetStream.Publish(subject, data)
	if err != nil {
		log.Printf("Failed to publish event to subject %s: %v", subject, err)
		return nil, err
	}

	return ack, nil
}
