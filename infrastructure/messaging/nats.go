package messaging

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/neodata-io/neodata-go/config"
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
	conn, err := nats.Connect(cfg.Messaging.PubsubBroker)
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

// CreateStreams sets up multiple JetStream streams based on the configuration
func (client *NATSClient) CreateStreams(cfg *config.AppConfig) error {
	// Iterate through the streams defined in the config
	for _, streamConfig := range cfg.Messaging.Streams {
		// Set the storage type (file or memory)
		storage := nats.FileStorage
		if streamConfig.StorageType == "memory" {
			storage = nats.MemoryStorage
		} else {
			storage = nats.FileStorage // Default to disk-based storage
		}

		// Define the stream configuration
		natsStreamConfig := &nats.StreamConfig{
			Name:     streamConfig.StreamName,
			Subjects: streamConfig.Subjects,
			MaxAge:   streamConfig.MaxAge,
			Storage:  storage,
			Replicas: streamConfig.Replicas,
		}

		// Create the stream
		_, err := client.JetStream.AddStream(natsStreamConfig)
		if err != nil {
			return fmt.Errorf("failed to create JetStream stream %s: %v", streamConfig.StreamName, err)
		}

		log.Printf("Stream %s created successfully with subjects: %v", streamConfig.StreamName, streamConfig.Subjects)
	}

	return nil
}
