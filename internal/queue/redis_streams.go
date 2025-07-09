package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sammyqtran/url-shortener/internal/events"
	"go.uber.org/zap"
)

// RedisStreamsQueue implements MessageQueue using Redis Streams
type RedisStreamsQueue struct {
	client *redis.Client
	config *StreamConfig
	logger *zap.Logger
}

// constructor for a new message queue using Redis Streams
func NewRedisStreamsQueue(client *redis.Client, config *StreamConfig, logger *zap.Logger) *RedisStreamsQueue {

	if config == nil {
		config = DefaultStreamConfig()
	}

	return &RedisStreamsQueue{
		client: client,
		config: config,
		logger: logger,
	}
}

// Publish sends an event to a specified redis stream
func (r *RedisStreamsQueue) Publish(ctx context.Context, stream string, event interface{}) error {

	// map to hold our event
	var eventMap map[string]interface{}

	switch e := event.(type) {
	case events.URLCreatedEventData:
		eventMap = e.ToMap()
	case events.URLAccessedEventData:
		eventMap = e.ToMap()
	default:
		r.logger.Error("Unsupported event type", zap.String("eventType", fmt.Sprintf("%T", event)))
		return fmt.Errorf("unsupported event type: %T", event)
	}

	// add to redis stream
	args := &redis.XAddArgs{
		Stream: stream,
		Values: eventMap,
	}

	_, err := r.client.XAdd(ctx, args).Result()
	if err != nil {
		r.logger.Error("Failed to publish event to stream", zap.String("stream", stream), zap.Error(err))
		return fmt.Errorf("failed to publish event to stream %s: %w", stream, err)
	}

	r.logger.Info("Published event to stream", zap.String("stream", stream), zap.String("event_type", eventMap["event_type"].(string)))
	return nil
}

func (r *RedisStreamsQueue) Subscribe(ctx context.Context, stream string, consumerGroup string, consumer string, handler EventHandler) error {

	err := r.createConsumerGroup(ctx, stream, consumerGroup)

	if err != nil {
		r.logger.Error("Failed to create consumer group", zap.Error(err))
		return fmt.Errorf("failed to create consumer group: %w", err)
	}

	r.logger.Info("Starting consumer in group for stream", zap.String("consumer", consumer), zap.String("consumerGroup", consumerGroup), zap.String("stream", stream))

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("Consumer stopping due to context cancellation", zap.String("consumer", consumer))
			return ctx.Err()

		default:
			// read from stream
			streams, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    consumerGroup,
				Consumer: consumer,
				Streams:  []string{stream, ">"},
				Count:    1,
				Block:    time.Duration(r.config.BlockTime) * time.Millisecond,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue
				}
				r.logger.Error("Error reading from stream", zap.Error(err))
				time.Sleep(5 * time.Second)
				continue
			}
			// Process messages
			for _, streamData := range streams {
				for _, message := range streamData.Messages {
					err := r.processMessage(ctx, message, handler)
					if err != nil {
						r.logger.Error("Error processing message", zap.String("messageID", message.ID), zap.Error(err))
						continue
					}

					// Acknowledge message
					err = r.client.XAck(ctx, stream, consumerGroup, message.ID).Err()
					if err != nil {
						r.logger.Error("Error acknowledging message", zap.String("messageID", message.ID), zap.Error(err))
					}
				}
			}

		}

	}
}

// createConsumerGroup creates a consumer group for the stream
func (r *RedisStreamsQueue) createConsumerGroup(ctx context.Context, stream, group string) error {
	// Try to create the consumer group
	err := r.client.XGroupCreateMkStream(ctx, stream, group, "0").Err()
	if err != nil {
		// If group already exists, that's fine
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			r.logger.Info("Consumer Group name already exists", zap.String("group", group))
			return nil
		}

		// If stream doesn't exist, create it first
		if err.Error() == "ERR The XGROUP subcommand requires the key to exist" {
			// Create stream with a dummy message
			_, err = r.client.XAdd(ctx, &redis.XAddArgs{
				Stream: stream,
				Values: map[string]interface{}{
					"init": "true",
				},
			}).Result()
			if err != nil {
				r.logger.Error("Failed to initalize stream", zap.Error(err))
				return fmt.Errorf("failed to initialize stream: %w", err)
			}

			// Now create the consumer group
			err = r.client.XGroupCreate(ctx, stream, group, "0").Err()
			if err != nil {
				r.logger.Error("Failed to create consumer group after stream init", zap.Error(err))
				return fmt.Errorf("failed to create consumer group after stream init: %w", err)
			}
		} else {
			r.logger.Error("Failed to create consumer group", zap.Error(err))
			return fmt.Errorf("failed to create consumer group: %w", err)
		}
	}

	return nil
}

// processMessage processes a single message from the stream
func (r *RedisStreamsQueue) processMessage(ctx context.Context, message redis.XMessage, handler EventHandler) error {
	// Extract event type
	eventTypeStr, ok := message.Values["event_type"].(string)
	if !ok {
		r.logger.Error("Missing or invalid event_type in message", zap.String("messageID", message.ID))
		return fmt.Errorf("missing or invalid event_type in message")
	}

	eventType := events.EventType(eventTypeStr)

	// Extract event data
	dataStr, ok := message.Values["data"].(string)
	if !ok {
		r.logger.Error("missing or invalid data in message")
		return fmt.Errorf("missing or invalid data in message")
	}

	// Call the handler
	return handler(ctx, eventType, []byte(dataStr))
}

// Close closes the Redis connection
func (r *RedisStreamsQueue) Close() error {
	return r.client.Close()
}

// Publisher wraps the message queue for easy publishing
type Publisher struct {
	queue  MessageQueue
	stream string
}

// NewPublisher creates a new event publisher
func NewPublisher(queue MessageQueue, stream string) *Publisher {
	return &Publisher{
		queue:  queue,
		stream: stream,
	}
}

// PublishURLCreated publishes a URL created event
func (p *Publisher) PublishURLCreated(ctx context.Context, shortCode, originalURL, createdBy string) error {
	event := events.URLCreatedEventData{
		BaseEvent: events.BaseEvent{
			ID:        generateEventID(),
			Type:      events.URLCreatedEvent,
			Timestamp: time.Now(),
			Source:    "gateway-service",
		},
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedBy:   createdBy,
	}

	return p.queue.Publish(ctx, p.stream, event)
}

// PublishURLAccessed publishes a URL accessed event
func (p *Publisher) PublishURLAccessed(ctx context.Context, shortCode, originalURL, userAgent, ipAddress, referrer string) error {
	event := events.URLAccessedEventData{
		BaseEvent: events.BaseEvent{
			ID:        generateEventID(),
			Type:      events.URLAccessedEvent,
			Timestamp: time.Now(),
			Source:    "gateway-service",
		},
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		UserAgent:   userAgent,
		IPAddress:   ipAddress,
		Referrer:    referrer,
	}

	return p.queue.Publish(ctx, p.stream, event)
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}
