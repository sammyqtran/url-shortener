package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sammyqtran/url-shortener/internal/analytics"
	"github.com/sammyqtran/url-shortener/internal/queue"
)

func main() {
	// Configuration
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0

	// Connect to Redis
	log.Printf("Connecting to Redis at %s", redisAddr)
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
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis")

	// Setup message queue
	streamConfig := queue.DefaultStreamConfig()
	messageQueue := queue.NewRedisStreamsQueue(redisClient, streamConfig)

	// Create analytics service
	analyticsService := analytics.NewAnalyticsService(messageQueue)

	// Create context for graceful shutdown
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	// Start analytics service in a goroutine
	go func() {
		if err := analyticsService.Start(ctx); err != nil {
			log.Printf("Analytics service error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down analytics service...")

	// Cancel context to stop the service
	cancel()

	// Give some time for graceful shutdown
	time.Sleep(2 * time.Second)

	// Close message queue
	if err := messageQueue.Close(); err != nil {
		log.Printf("Error closing message queue: %v", err)
	}

	log.Println("Analytics service stopped")
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
