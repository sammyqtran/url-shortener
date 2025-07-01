package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/sammyqtran/url-shortener/internal/gateway"
	"github.com/sammyqtran/url-shortener/internal/queue"
	pb "github.com/sammyqtran/url-shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0

	conn, err := grpc.NewClient("url-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Connect to Redis for message queue
	log.Printf("Connecting to Redis at %s", redisAddr)
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
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis")

	// Setup message queue
	streamConfig := queue.DefaultStreamConfig()
	messageQueue := queue.NewRedisStreamsQueue(redisClient, streamConfig)
	publisher := queue.NewPublisher(messageQueue, streamConfig.URLEventsStream)

	grpcClient := pb.NewURLServiceClient(conn)

	server := &gateway.GatewayServer{
		GrpcClient: grpcClient,
		Publisher:  publisher,
	}

	r := mux.NewRouter()

	r.HandleFunc("/create", server.HandleCreateShortURL).Methods("POST")
	r.HandleFunc("/healthz", server.HandleHealthCheck).Methods("GET")
	r.HandleFunc("/{shortCode}", server.HandleGetOriginalURL).Methods("GET")

	log.Println("Gateway service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
