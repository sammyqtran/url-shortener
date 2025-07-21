package analytics

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/sammyqtran/url-shortener/internal/events"
	"github.com/sammyqtran/url-shortener/internal/metrics"
	"github.com/sammyqtran/url-shortener/internal/queue"
	"go.uber.org/zap"
)

type MockMessageQueue struct {
	PublishedEvents []interface{}
	Subscribed      bool
	Handler         queue.EventHandler
}

// sends an event to the specificed stream/topic
func (m *MockMessageQueue) Publish(ctx context.Context, stream string, event interface{}) error {
	m.PublishedEvents = append(m.PublishedEvents, event)
	return nil
}

// Subscribe starts consuming events from the specified stream/topic
func (m *MockMessageQueue) Subscribe(ctx context.Context, stream string, consumerGroup string, consumer string, handler queue.EventHandler) error {
	m.Subscribed = true
	m.Handler = handler
	return nil
}

func (m *MockMessageQueue) TriggerEvent(ctx context.Context, eventType events.EventType, data []byte) error {
	if m.Handler != nil {
		return m.Handler(ctx, eventType, data)
	}
	return nil
}

// closes the message queue connection
func (m *MockMessageQueue) Close() error {
	return nil
}

func TestHandleURLCreated(t *testing.T) {
	mockQueue := new(MockMessageQueue)
	mockMetrics := &metrics.NoopMetrics{}
	service := &AnalyticsService{
		MessageQueue: mockQueue,
		Logger:       zap.NewNop(),
		Metrics:      mockMetrics,
	}
	event := events.URLCreatedEventData{
		BaseEvent: events.BaseEvent{
			ID:        "uuid-fake",
			Type:      events.URLCreatedEvent,
			Timestamp: time.Now(),
			Source:    "gateway-service",
		},
		ShortCode:   "abc123",
		OriginalURL: "https://google.com",
		CreatedBy:   "great-ruler",
	}
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}

	err = service.handleURLCreated(context.Background(), data)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	invalidJSON := []byte(`{
  "short_code": 12345,
  "original_url": "https://google.com",
  "created_by": true
	}`)
	err = service.handleURLCreated(context.Background(), invalidJSON)

	if err == nil {
		t.Errorf("expected error, got none")
	}

}

func TestHandleURLAccessed(t *testing.T) {
	mockMetrics := &metrics.NoopMetrics{}
	mockQueue := new(MockMessageQueue)
	service := &AnalyticsService{
		MessageQueue: mockQueue,
		Logger:       zap.NewNop(),
		Metrics:      mockMetrics,
	}
	event := events.URLAccessedEventData{
		BaseEvent: events.BaseEvent{
			ID:        "fake-uuid",
			Type:      events.URLAccessedEvent,
			Timestamp: time.Now(),
			Source:    "gateway-service",
		},
		ShortCode:   "abc123",
		OriginalURL: "https://google.com",
		UserAgent:   "lord-ruler",
		IPAddress:   "192.168.1.101",
		Referrer:    "hamburger-dude",
	}
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}
	err = service.handleURLAccessed(context.Background(), data)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	invalidJSON := []byte(`{
  "short_code": 12345,
  "original_url": "https://google.com",
  "created_by": true
	}`)
	err = service.handleURLAccessed(context.Background(), invalidJSON)

	if err == nil {
		t.Errorf("expected error, got none")
	}

}

func TestHandleEvent(t *testing.T) {
	mockMetrics := &metrics.NoopMetrics{}
	mockQueue := new(MockMessageQueue)
	service := &AnalyticsService{
		MessageQueue: mockQueue,
		Logger:       zap.NewNop(),
		Metrics:      mockMetrics,
	}
	event := events.URLCreatedEventData{
		BaseEvent: events.BaseEvent{
			ID:        "uuid-fake",
			Type:      events.URLCreatedEvent,
			Timestamp: time.Now(),
			Source:    "gateway-service",
		},
		ShortCode:   "abc123",
		OriginalURL: "https://google.com",
		CreatedBy:   "great-ruler",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}
	err = service.handleEvent(context.Background(), events.URLCreatedEvent, data)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	event2 := events.URLAccessedEventData{
		BaseEvent: events.BaseEvent{
			ID:        "fake-uuid",
			Type:      events.URLAccessedEvent,
			Timestamp: time.Now(),
			Source:    "gateway-service",
		},
		ShortCode:   "abc123",
		OriginalURL: "https://google.com",
		UserAgent:   "lord-ruler",
		IPAddress:   "192.168.1.101",
		Referrer:    "hamburger-dude",
	}
	data, err = json.Marshal(event2)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}
	err = service.handleEvent(context.Background(), events.URLAccessedEvent, data)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	err = service.handleEvent(context.Background(), "non.existent", data)
	if err != nil {
		t.Errorf("expected no error for unknown event, got: %v", err)
	}
}

func TestStart(t *testing.T) {
	mockMetrics := &metrics.NoopMetrics{}
	mockQueue := new(MockMessageQueue)
	service := &AnalyticsService{
		MessageQueue: mockQueue,
		Logger:       zap.NewNop(),
		Metrics:      mockMetrics,
	}

	err := service.Start(context.Background())
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if mockQueue.Subscribed != true {
		t.Errorf("expected Subscribe to be called")
	}
}
