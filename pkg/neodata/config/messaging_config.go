// config/messaging_config.go
package config

import "time"

// MessagingConfig holds the messaging system configuration
type MessagingConfig struct {
	PubsubBackend string             `mapstructure:"pubsub_backend"`
	PubsubBroker  string             `mapstructure:"pubsub_broker"`
	Streams       []NATSStreamConfig `mapstructure:"streams"`
}

// NATSStreamConfig defines configuration for a single JetStream stream
type NATSStreamConfig struct {
	StreamName  string        `mapstructure:"stream_name"`
	Subjects    []string      `mapstructure:"subjects"`
	MaxAge      time.Duration `mapstructure:"max_age"`
	StorageType string        `mapstructure:"storage_type"`
	Replicas    int           `mapstructure:"replicas"`
}
