package messaging

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

// Messaging defines the interface for publishing messages
type Messaging interface {
	Publish(ctx context.Context, subject string, data []byte) (*jetstream.PubAck, error)
}
