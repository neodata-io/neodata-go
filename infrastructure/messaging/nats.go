package messaging

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type NATSConfig struct {
	URL     string
	Timeout time.Duration
}

type NATSClient struct {
	Conn      *nats.Conn
	JetStream nats.JetStreamContext
}

// NewNATSClient creates a new NATS JetStream connection
func NewNATSClient(cfg *config.AppConfig) (*NATSClient, error) {
	conn, err := nats.Connect(cfg.Messaging.URL, nats.Timeout(config.Timeout))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
		return nil, err
	}

	js, err := conn.JetStream()
	if err != nil {
		log.Fatalf("Failed to get JetStream context: %v", err)
		return nil, err
	}

	return &NATSClient{Conn: conn, JetStream: js}, nil
}
