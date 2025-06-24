package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
}

func NewURLService(repo repository.URLRepository) *URLService {
	service := &URLService{
		repo:    repo,
		baseURL: "http://localhost:8080/",
	}
	service.codeGenerator = service.GenerateShortCode
	return service
}

func (s *URLService) CreateShortURL(ctx context.Context, req *pb.CreateURLRequest) (*pb.CreateURLResponse, error) {

	// url validation
	if err := s.validateURL(req.OriginalUrl); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid URL: %v", err)
	}

	// short code generation
	shortCode, err := s.codeGenerator(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate short code: %v", err)
	}

	// Create URL model
	urlModel := &models.URL{
		UserID:      req.UserId,
		ShortCode:   shortCode,
		OriginalURL: req.OriginalUrl,
		// ExpiresAt:   expiresAt,
	}

	// Save to database
	if err := s.repo.Create(ctx, urlModel); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create URL: %v", err)
	}

	return &pb.CreateURLResponse{
		ShortCode: shortCode,
		ShortUrl:  s.baseURL + shortCode, // Replace with actual base URL
		Success:   true,
		Error:     "",
	}, nil
}

func (s *URLService) GetOriginalURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	if req.ShortCode == "" {
		return nil, status.Error(codes.InvalidArgument, "short_code cannot be empty")
	}

	// Retrieve from repository
	urlModel, err := s.repo.GetByShortCode(ctx, req.ShortCode)
	if err != nil {
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

	// increment click count asynchronously
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.repo.IncrementClickCount(ctx, req.ShortCode); err != nil {
			// Log the error, but do not return it to the user
			fmt.Printf("failed to increment click count for %s: %v\n", req.ShortCode, err)
		}
	}()

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
