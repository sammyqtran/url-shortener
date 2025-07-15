package analytics

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sammyqtran/url-shortener/internal/events"
	"github.com/sammyqtran/url-shortener/internal/metrics"
	"github.com/sammyqtran/url-shortener/internal/queue"
	"go.uber.org/zap"
)

type AnalyticsService struct {
	MessageQueue queue.MessageQueue
	Logger       *zap.Logger
	Metrics      *metrics.PrometheusMetrics
	// Add database connection here when you want to persist analytics
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(messageQueue queue.MessageQueue, Logger *zap.Logger, metrics *metrics.PrometheusMetrics) *AnalyticsService {
	return &AnalyticsService{
		MessageQueue: messageQueue,
		Logger:       Logger,
		Metrics:      metrics,
	}
}

// Start starts the analytics service
func (a *AnalyticsService) Start(ctx context.Context) error {
	a.Logger.Info("Starting analytics service...")

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

	service := "analytics-service"
	eventTimer := time.Now()

	defer func() {
		a.Metrics.IncConsumeEvent(service, string(eventType))
		a.Metrics.ObserveConsumeEventLatency(service, string(eventType), time.Since(eventTimer).Seconds())
	}()

	a.Logger.Info("Received event", zap.String("eventType", string(eventType)))
	switch eventType {
	case events.URLCreatedEvent:
		return a.handleURLCreated(ctx, data)
	case events.URLAccessedEvent:
		return a.handleURLAccessed(ctx, data)
	default:
		a.Logger.Warn("Unknown event type", zap.String("event_type", string(eventType)))
		return nil
	}
}

// handleURLAccessed processes URL access events
func (a *AnalyticsService) handleURLAccessed(ctx context.Context, data []byte) error {
	var event events.URLAccessedEventData
	if err := json.Unmarshal(data, &event); err != nil {
		a.Metrics.IncConsumeEventError("analytics-service", string(events.URLAccessedEvent))
		return err
	}

	a.Logger.Info("URL Accessed",
		zap.String("shortCode", event.ShortCode),
		zap.String("originalURL", event.OriginalURL),
		zap.String("userAgent", event.UserAgent),
		zap.String("ip", event.IPAddress),
		zap.String("referrer", event.Referrer),
		zap.String("timestamp", event.Timestamp.Format(time.RFC3339)),
	)

	// TODO: Store analytics data in database
	// Example: Insert into url_accesses table with all the tracking data

	return nil
}

func (a *AnalyticsService) handleURLCreated(ctx context.Context, data []byte) error {
	var event events.URLCreatedEventData
	if err := json.Unmarshal(data, &event); err != nil {
		a.Metrics.IncConsumeEventError("analytics-service", string(events.URLCreatedEvent))
		return err
	}
	a.Logger.Info("URL Created",
		zap.String("shortCode", event.ShortCode),
		zap.String("originalURL", event.OriginalURL),
		zap.String("createdBy", event.CreatedBy),
		zap.String("timestamp", event.Timestamp.Format(time.RFC3339)),
	)

	return nil
}
