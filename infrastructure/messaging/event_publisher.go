package messaging

import (
	"log"

	"github.com/nats-io/nats.go"
)

// EventPublisher handles publishing events to NATS JetStream
type EventPublisher struct {
	jetStream nats.JetStreamContext
}

// NewEventPublisher creates a new instance of EventPublisher
func NewEventPublisher(client *NATSClient) *EventPublisher {
	return &EventPublisher{jetStream: client.JetStream}
}

// Publish publishes an event to a specific subject
func (p *EventPublisher) Publish(subject string, data []byte) error {
	_, err := p.jetStream.Publish(subject, data)
	if err != nil {
		log.Printf("Failed to publish event to subject %s: %v", subject, err)
		return err
	}
	log.Printf("Event published to subject: %s", subject)
	return nil
}
