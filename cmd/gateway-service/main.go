package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/sammyqtran/url-shortener/internal/gateway"
	"github.com/sammyqtran/url-shortener/internal/queue"
	pb "github.com/sammyqtran/url-shortener/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// Structured Logging
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	redisAddr := getEnv("REDIS_ADDR", "redis:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0

	target := getEnv("URL_SERVICE_HOST", "url-service")
	conn, err := grpc.NewClient(target+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Error connecting to Redis", zap.Error(err))
	}
	defer conn.Close()

	// Connect to Redis for message queue
	logger.Info("Conntecting to Redis at", zap.String("redisAddr", redisAddr))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Fatal("Failed to Connect to Redis:", zap.Error(err))
	}
	logger.Info("Successfully connected to Redis")

	// Setup message queue
	streamConfig := queue.DefaultStreamConfig()
	messageQueue := queue.NewRedisStreamsQueue(redisClient, streamConfig, logger)
	publisher := queue.NewPublisher(messageQueue, streamConfig.URLEventsStream)

	grpcClient := pb.NewURLServiceClient(conn)

	server := &gateway.GatewayServer{
		GrpcClient: grpcClient,
		Publisher:  publisher,
		Logger:     logger,
	}

	r := mux.NewRouter()

	r.HandleFunc("/create", server.HandleCreateShortURL).Methods("POST")
	r.HandleFunc("/healthz", server.HandleHealthCheck).Methods("GET")
	r.HandleFunc("/{shortCode}", server.HandleGetOriginalURL).Methods("GET")

	logger.Info("Gateway service listening on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Fatal("HTTP server failed", zap.Error(err))
	}

}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
