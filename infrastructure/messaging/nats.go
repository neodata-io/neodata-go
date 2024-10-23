package messaging

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/neodata-io/neodata-go/config"
)

type NATSClient struct {
	nc *nats.Conn
	js jetstream.JetStream
}

// NewNATSClient creates a new NATS JetStream connection
func NewNATSClient(ctx context.Context, natsURL string) (*NATSClient, error) {
	// In the `jetstream` package, almost all API calls rely on `context.Context` for timeout/cancellation handling
	_, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create a JetStream management interface  to manage and list streams
	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	return &NATSClient{nc: nc, js: js}, nil
}

// Close closes the NATS connection
func (n *NATSClient) Close() {
	if n.nc != nil {
		n.nc.Close()
	}
}

// CreateStreams sets up multiple JetStream streams based on the configuration
func (n *NATSClient) CreateStreams(ctx context.Context, cfg *config.AppConfig) error {
	// Iterate through the streams defined in the config
	for _, streamConfig := range cfg.Messaging.Streams {
		// Set the storage type (file or memory)
		storage := jetstream.FileStorage
		if streamConfig.StorageType == "memory" {
			storage = jetstream.MemoryStorage
		}

		// create a stream (this is an idempotent operation)
		_, err := n.js.CreateStream(ctx, jetstream.StreamConfig{
			Name:     streamConfig.StreamName,
			Subjects: streamConfig.Subjects,
			MaxAge:   streamConfig.MaxAge,
			Storage:  storage,
			Replicas: streamConfig.Replicas,
		})

		if err != nil {
			return fmt.Errorf("failed to create JetStream stream %s: %v", streamConfig.StreamName, err)
		}

		log.Printf("Stream %s created successfully with subjects: %v", streamConfig.StreamName, streamConfig.Subjects)
	}

	return nil
}

// Create durable consumer
func (n *NATSClient) CreateConsumers(ctx context.Context, cfg *config.AppConfig) error {
	// Iterate through the streams defined in the config
	for _, streamConfig := range cfg.Messaging.Streams {
		// Create durable consumer
		_, err := n.js.CreateOrUpdateConsumer(ctx, streamConfig.StreamName, jetstream.ConsumerConfig{
			Durable:   streamConfig.StreamName + "_consumer",
			AckPolicy: jetstream.AckExplicitPolicy,
		})

		if err != nil {
			return fmt.Errorf("failed to create Consumer %s: %v", streamConfig.StreamName, err)
		}
	}
	return nil
}
