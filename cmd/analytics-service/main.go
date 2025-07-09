package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sammyqtran/url-shortener/internal/analytics"
	"github.com/sammyqtran/url-shortener/internal/queue"
	"go.uber.org/zap"
)

func main() {

	// Structured Logging setup
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Configuration
	redisAddr := getEnv("REDIS_ADDR", "redis:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0

	// Connect to Redis
	logger.Info("Connecting to Redis", zap.String("redis_addr", redisAddr))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Successfully connected to Redis")

	// Setup message queue
	streamConfig := queue.DefaultStreamConfig()
	messageQueue := queue.NewRedisStreamsQueue(redisClient, streamConfig, logger)

	// Create analytics service
	analyticsService := analytics.NewAnalyticsService(messageQueue, logger)

	// Create context for graceful shutdown
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Start analytics service in a goroutine
	go func() {
		if err := analyticsService.Start(ctx); err != nil {
			logger.Error("Analytics service error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down analytics service...")

	// Cancel context to stop the service
	cancel()

	// Give some time for graceful shutdown
	time.Sleep(2 * time.Second)

	// Close message queue
	if err := messageQueue.Close(); err != nil {
		logger.Error("Error closing message queue", zap.Error(err))
	}

	logger.Info("Analytics service stopped")
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
