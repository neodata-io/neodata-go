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
func NewNATSClient(natsURL string) (*NATSClient, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	return &NATSClient{Conn: conn, JetStream: js}, nil
}

func (n *NATSClient) Close() {
	if n.Conn != nil {
		n.Conn.Close()
	}
}

// CreateStreams sets up multiple JetStream streams based on the configuration
func (client *NATSClient) CreateStreams(cfg *config.AppConfig) error {
	// Iterate through the streams defined in the config
	for _, streamConfig := range cfg.Messaging.Streams {
		// Set the storage type (file or memory)
		storage := nats.FileStorage
		if streamConfig.StorageType == "memory" {
			storage = nats.MemoryStorage
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
