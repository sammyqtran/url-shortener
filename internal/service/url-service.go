package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/redis/go-redis/v9"
	"github.com/sammyqtran/url-shortener/internal/metrics"
	"github.com/sammyqtran/url-shortener/internal/models"
	"github.com/sammyqtran/url-shortener/internal/repository"
	pb "github.com/sammyqtran/url-shortener/proto"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLService struct {
	pb.UnimplementedURLServiceServer
	repo          repository.URLRepository
	baseURL       string
	codeGenerator func(ctx context.Context) (string, error)
	cache         *redis.Client
	Logger        *zap.Logger
	Metrics       *metrics.PrometheusMetrics
}

func NewURLService(repo repository.URLRepository, cache *redis.Client, logger *zap.Logger, metrics *metrics.PrometheusMetrics) *URLService {
	service := &URLService{
		repo:    repo,
		baseURL: "http://localhost:8080/",
		cache:   cache,
		Logger:  logger,
		Metrics: metrics,
	}
	service.codeGenerator = service.GenerateShortCode
	return service
}

func (s *URLService) CreateShortURL(ctx context.Context, req *pb.CreateURLRequest) (*pb.CreateURLResponse, error) {
	service := "url-service"

	// url validation
	if err := s.validateURL(req.OriginalUrl); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid URL: %v", err)
	}

	// short code generation
	shortCode, err := s.codeGenerator(ctx)
	if err != nil {
		s.Logger.Error("Failed to generate short code", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to generate short code: %v", err)
	}

	now := time.Now()

	// Create URL model
	urlModel := &models.URL{
		UserID:      req.UserId,
		ShortCode:   shortCode,
		OriginalURL: req.OriginalUrl,
		CreatedAt:   now,
		UpdatedAt:   now,
		ClickCount:  0,
		// ExpiresAt:   expiresAt,
	}

	// record db operation and duration
	s.Metrics.IncDBOperation(service, "Create")
	dbTimer := time.Now()
	// Save to database
	if err := s.repo.Create(ctx, urlModel); err != nil {
		s.Metrics.IncDBError(service, "Create")
		s.Logger.Error("Failed to create short URL", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create URL: %v", err)
	}
	s.Metrics.ObserveDBOperationDuration(service, "Create", time.Since(dbTimer).Seconds())

	// Put in cache

	err = s.setCacheFromModel(ctx, shortCode, urlModel)
	if err != nil {
		s.Metrics.IncCacheError(service, "map_url", "set")
		s.Logger.Error("Failed to cache short URL", zap.Error(err))
	}

	return &pb.CreateURLResponse{
		ShortCode: shortCode,
		ShortUrl:  s.baseURL + shortCode,
		Success:   true,
		Error:     "",
	}, nil
}

func (s *URLService) GetOriginalURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	if req.ShortCode == "" {
		return nil, status.Error(codes.InvalidArgument, "short_code cannot be empty")
	}

	// Try cache first
	cachedURL, err := s.getFromCache(ctx, req.ShortCode)
	if err == nil {
		s.Metrics.IncCacheHit("url-service", "map_url")
		s.Logger.Info("Cache hit", zap.String("shortCode", req.ShortCode))
		// Cache hit - check expiration
		if cachedURL.ExpiresAt != nil && cachedURL.ExpiresAt.Before(time.Now()) {
			// URL expired - remove from cache
			s.removeFromCache(ctx, req.ShortCode)
			return &pb.GetURLResponse{
				Found: false,
				Error: "URL has expired",
			}, nil
		}

		// Increment click count asynchronously (don't block response)
		go s.incrementClickCountAsync(req.ShortCode)

		return &pb.GetURLResponse{
			OriginalUrl: cachedURL.OriginalURL,
			Found:       true,
		}, nil
	}
	s.Logger.Info("Cache miss", zap.String("shortCode", req.ShortCode))

	// Fall back retrieve from repository
	s.Metrics.IncDBOperation("url-service", "GetByShortCode")
	dbTimer := time.Now()
	urlModel, err := s.repo.GetByShortCode(ctx, req.ShortCode)
	s.Metrics.ObserveDBOperationDuration("url-service", "GetByShortCode", time.Since(dbTimer).Seconds())
	if err != nil {
		s.Metrics.IncDBError("url-service", "GetByShortCode")
		s.Logger.Error("Error retrieving from repository", zap.Error(err))
		if err == repository.ErrURLNotFound {
			return &pb.GetURLResponse{
				Found: false,
				Error: "URL not found",
			}, nil
		}
		return &pb.GetURLResponse{
			Found: false,
			Error: fmt.Sprintf("failed to retrieve URL: %v", err),
		}, nil
	}
	// Check if the URL has expired
	if urlModel.ExpiresAt != nil && urlModel.ExpiresAt.Before(time.Now()) {
		return &pb.GetURLResponse{
			Found: false,
			Error: "URL has expired",
		}, nil
	}
	s.Logger.Info("Database fetch", zap.String("shortCode", req.ShortCode))

	// populate cache from db
	s.setCacheFromModel(ctx, req.ShortCode, urlModel)

	// increment click count asynchronously
	go s.incrementClickCountAsync(req.ShortCode)

	// Return the original URL if found
	return &pb.GetURLResponse{
		OriginalUrl: urlModel.OriginalURL,
		Found:       true,
	}, nil
}

func (s *URLService) HealthCheck(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Healthy: true,
	}, nil
}

// validateURL checks if the provided URL is valid.
func (s *URLService) validateURL(urlStr string) error {
	if urlStr == "" {
		return status.Error(codes.InvalidArgument, "URL cannot be empty")
	}
	// Optionally, add more robust URL validation here.
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must have a scheme (http or https)")
	}

	if !strings.Contains(parsedURL.Host, ".") {
		return fmt.Errorf("URL must contain a valid domain")
	}
	return nil
}

// Temporary short code generator - will be improved later
func generateRandomCode() string {

	const length = 6

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *URLService) GenerateShortCode(ctx context.Context) (string, error) {
	// In a real application, you would check for uniqueness in the database.
	// Here we just generate a random code.
	const maxAttempts = 10
	for attempts := 0; attempts < maxAttempts; attempts++ {
		shortCode := generateRandomCode()
		exists, err := s.repo.IsShortCodeExists(ctx, shortCode)
		if err != nil {
			return "", fmt.Errorf("error checking short code existence: %w", err)
		}

		if !exists {
			return shortCode, nil
		}
	}

	return "", fmt.Errorf("failed to generate a unique short code after %d attempts", maxAttempts)
}

// async function to increment click count

func (s *URLService) incrementClickCountAsync(shortCode string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.repo.IncrementClickCount(ctx, shortCode); err != nil {
		s.Logger.Error("Failed to increment click count", zap.String("shortCode", shortCode), zap.Error(err))
	}
}

// Redis helpers to cache URL model

func (s *URLService) setCacheFromModel(ctx context.Context, shortCode string, urlModel *models.URL) error {
	cacheKey := fmt.Sprintf("url:%s", shortCode)

	data, err := json.Marshal(urlModel)
	if err != nil {
		s.Logger.Error("Failed to marshal cache data", zap.String("shortCode", shortCode), zap.Error(err))
		return fmt.Errorf("failed to marshal cache data for %s: %w", shortCode, err)
	}

	// Set cache with TTL
	err = s.cache.Set(ctx, cacheKey, data, 10*time.Minute).Err()
	if err != nil {
		s.Logger.Error("Failed to set cache", zap.String("shortCode", shortCode), zap.Error(err))
		return fmt.Errorf("failed to set cache for %s: %w", shortCode, err)
	}

	return nil
}

func (s *URLService) getFromCache(ctx context.Context, shortCode string) (*models.URL, error) {
	cacheKey := fmt.Sprintf("url:%s", shortCode)

	result, err := s.cache.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		s.Metrics.IncCacheMiss("url-service", "map_url")
		return nil, fmt.Errorf("cache miss")
	}
	if err != nil {
		s.Metrics.IncCacheError("url-service", "map_url", "get")
		// Log cache error but don't fail the request
		s.Logger.Warn("Cache get error", zap.String("shortCode", shortCode), zap.Error(err))
		return nil, fmt.Errorf("cache error")
	}

	var urlModel models.URL
	err = json.Unmarshal([]byte(result), &urlModel)
	if err != nil {
		s.Logger.Error("Failed to unmarshal cached URL", zap.String("shortCode", shortCode), zap.Error(err))
		return nil, fmt.Errorf("cache unmarshal error")
	}

	return &urlModel, nil
}
func (s *URLService) removeFromCache(ctx context.Context, shortCode string) {
	cacheKey := fmt.Sprintf("url:%s", shortCode)
	err := s.cache.Del(ctx, cacheKey).Err()
	if err != nil {
		s.Metrics.IncCacheError("url-service", "map_url", "delete")
		s.Logger.Error("Failed to remove from cache", zap.String("shortCode", shortCode), zap.Error(err))
	}
}
