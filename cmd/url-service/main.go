package main

import (
	"context"
	"net"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/redis/go-redis/v9"
	"github.com/sammyqtran/url-shortener/internal/database"
	"github.com/sammyqtran/url-shortener/internal/repository/postgres"
	"github.com/sammyqtran/url-shortener/internal/service"
	pb "github.com/sammyqtran/url-shortener/proto"
)

func main() {
	// Structured Logging
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// database configuration from environment variables
	dbConfig := database.Config{
		Host:         getEnv("DB_HOST", "postgres"),
		Port:         getEnvAsInt("DB_PORT", 5432),
		User:         getEnv("DB_USER", "postgres"),
		Password:     getEnv("DB_PASSWORD", "password"),
		DatabaseName: getEnv("DB_NAME", "urlshortener"),
		SSLMode:      getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		MaxLifetime:  time.Duration(getEnvAsInt("DB_MAX_LIFETIME_MINUTES", 30)) * time.Minute,
	}
	// connect to the database
	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	logger.Info("Connected to database successfully")

	// run migrations
	if err := database.RunMigrations(db); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	logger.Info("Database migrations applied successfully")

	// create a new URL repository instance
	urlRepo := postgres.NewPostgresURLRepository(db, logger)

	// Create a Redis client and connect to Redis
	cache := redis.NewClient(&redis.Options{
		Addr: getEnv("REDIS_ADDR", "redis:6379"),
		DB:   0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := cache.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Connected to Redis successfully")

	// create service instance (uses default baseURL from service package)
	urlService := service.NewURLService(urlRepo, cache, logger)

	// create a new gRPC server
	logger.Info("Starting gRPC server on port 50051...")
	grpcServer := grpc.NewServer()

	pb.RegisterURLServiceServer(grpcServer, urlService)
	reflection.Register(grpcServer)

	// listen on port 50051
	port := getEnv("GRPC_PORT", "50051")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal("Failed to listen on port", zap.String("port", port), zap.Error(err))
	}
	logger.Info("gRPC server listening on port", zap.String("port", port))
	logger.Info("Database connection info",
		zap.String("user", dbConfig.User),
		zap.String("host", dbConfig.Host),
		zap.Int("port", dbConfig.Port),
		zap.String("database", dbConfig.DatabaseName),
	)
	// serve the gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatal("Failed to serve gRPC server", zap.Error(err))
	}
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
