// internal/events/types.go
package events

import (
	"encoding/json"
	"time"
)

// EventType represents the type of event being published
type EventType string

const (
	URLCreatedEvent  EventType = "url.created"
	URLAccessedEvent EventType = "url.accessed"
)

// BaseEvent contains common fields for all events
type BaseEvent struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

// URLCreatedEventData represents data for URL creation events
type URLCreatedEventData struct {
	BaseEvent
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	CreatedBy   string     `json:"created_by,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// URLAccessedEventData represents data for URL access events
type URLAccessedEventData struct {
	BaseEvent
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	UserAgent   string `json:"user_agent,omitempty"`
	IPAddress   string `json:"ip_address,omitempty"`
	Referrer    string `json:"referrer,omitempty"`
}

// ToJSON serializes the event to JSON
func (e BaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// ToMap converts event to map[string]interface{} for Redis Streams
func (e URLCreatedEventData) ToMap() map[string]interface{} {
	data, _ := json.Marshal(e)
	return map[string]interface{}{
		"event_type": string(e.Type),
		"data":       string(data),
		"timestamp":  e.Timestamp.Unix(),
	}
}

func (e URLAccessedEventData) ToMap() map[string]interface{} {
	data, _ := json.Marshal(e)
	return map[string]interface{}{
		"event_type": string(e.Type),
		"data":       string(data),
		"timestamp":  e.Timestamp.Unix(),
	}
}
