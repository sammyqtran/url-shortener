package analytics

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/sammyqtran/url-shortener/internal/events"
	"github.com/sammyqtran/url-shortener/internal/queue"
)

type AnalyticsService struct {
	MessageQueue queue.MessageQueue
	// Add database connection here when you want to persist analytics
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(messageQueue queue.MessageQueue) *AnalyticsService {
	return &AnalyticsService{
		MessageQueue: messageQueue,
	}
}

// Start starts the analytics service
func (a *AnalyticsService) Start(ctx context.Context) error {
	log.Println("Starting analytics service...")

	// Start consuming events
	streamConfig := queue.DefaultStreamConfig()
	return a.MessageQueue.Subscribe(
		ctx,
		streamConfig.URLEventsStream,
		streamConfig.ConsumerGroup,
		"analytics-consumer-1",
		a.handleEvent,
	)
}

// handleEvent processes incoming events
func (a *AnalyticsService) handleEvent(ctx context.Context, eventType events.EventType, data []byte) error {
	log.Printf("Received event: %s", eventType)

	switch eventType {
	case events.URLCreatedEvent:
		return a.handleURLCreated(ctx, data)
	case events.URLAccessedEvent:
		return a.handleURLAccessed(ctx, data)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}

// handleURLAccessed processes URL access events
func (a *AnalyticsService) handleURLAccessed(ctx context.Context, data []byte) error {
	var event events.URLAccessedEventData
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	log.Printf("URL Accessed - Short Code: %s, Original URL: %s, User Agent: %s, IP: %s, Referrer: %s, Timestamp: %s",
		event.ShortCode,
		event.OriginalURL,
		event.UserAgent,
		event.IPAddress,
		event.Referrer,
		event.Timestamp.Format(time.RFC3339),
	)

	// TODO: Store analytics data in database
	// Example: Insert into url_accesses table with all the tracking data

	return nil
}

func (a *AnalyticsService) handleURLCreated(ctx context.Context, data []byte) error {
	var event events.URLCreatedEventData
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	log.Printf("URL Created - Short Code: %s, Original URL: %s, Created By: %s, Timestamp: %s",
		event.ShortCode,
		event.OriginalURL,
		event.CreatedBy,
		event.Timestamp.Format(time.RFC3339),
	)

	return nil
}
