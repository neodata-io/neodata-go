package messaging

/* // EventHandler defines the interface for handling events
type EventHandler func(msg *nats.Msg) error

// EventSubscriber handles subscribing to events from NATS JetStream
type EventSubscriber struct {
	jetStream nats.JetStreamContext
}

// NewEventSubscriber creates a new instance of EventSubscriber
func NewEventSubscriber(client *NATSClient) *EventSubscriber {
	return &EventSubscriber{jetStream: client.JetStream}
}

// Subscribe subscribes to a specific subject and processes events with the provided handler
func (s *EventSubscriber) Subscribe(subject, durableName string, handler EventHandler) error {
	_, err := s.jetStream.Subscribe(subject, func(msg *nats.Msg) {
		if err := handler(msg); err != nil {
			log.Printf("Failed to process message: %v", err)
		}
	}, nats.Durable(durableName), nats.ManualAck())
	if err != nil {
		log.Printf("Failed to subscribe to subject %s: %v", subject, err)
		return err
	}
	log.Printf("Subscribed to subject: %s", subject)
	return nil
}
*/
