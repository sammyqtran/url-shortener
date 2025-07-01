package queue

import (
	"context"

	"github.com/sammyqtran/url-shortener/internal/events"
)

type MessageQueue interface {

	// sends an event to the specificed stream/topic
	Publish(ctx context.Context, stream string, event interface{}) error

	// Subscribe starts consuming events from the specified stream/topic
	Subscribe(ctx context.Context, stream string, consumerGroup string, consumer string, handler EventHandler) error

	// clolses the message queue connection
	Close() error
}

// defines the function signature for handling events
type EventHandler func(ctx context.Context, eventType events.EventType, data []byte) error

// StreamConfig holds configuration for Redis Streams
type StreamConfig struct {
	URLEventsStream string
	ConsumerGroup   string
	ReadTimeout     int // seconds
	BlockTime       int // milliseconds
}

// DefaultStreamConfig returns default configuration
func DefaultStreamConfig() *StreamConfig {
	return &StreamConfig{
		URLEventsStream: "url-events",
		ConsumerGroup:   "analytics-group",
		ReadTimeout:     30,
		BlockTime:       5000,
	}
}

// EventPublisher defines the methods needed to publish events.
// Implemented by Publisher and by mocks in tests.
type EventPublisher interface {
	PublishURLCreated(ctx context.Context, shortCode, originalURL, createdBy string) error
	PublishURLAccessed(ctx context.Context, shortCode, originalURL, userAgent, ipAddress, referrer string) error
}
