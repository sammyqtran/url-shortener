package service

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/sammyqtran/url-shortener/internal/models"
	pb "github.com/sammyqtran/url-shortener/proto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(ctx context.Context, url *models.URL) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *MockRepo) GetByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	args := m.Called(ctx, shortCode)

	return args.Get(0).(*models.URL), args.Error(1)
}

func (m *MockRepo) GetByID(ctx context.Context, id int64) (*models.URL, error) {
	return nil, nil
}

func (m *MockRepo) Update(ctx context.Context, url *models.URL) error {
	return nil
}

func (m *MockRepo) Delete(ctx context.Context, shortCode string) error {
	return nil
}

func (m *MockRepo) IncrementClickCount(ctx context.Context, shortCode string) error {
	return nil
}

// GetStats returns URL statistics
func (m *MockRepo) GetStats(ctx context.Context, shortCode string) (*models.URL, error) {
	return nil, nil
}

// ListURLs returns paginated list of URLs
func (m *MockRepo) ListURLs(ctx context.Context, limit, offset int) ([]*models.URL, error) {
	return nil, nil
}

// IsShortCodeExists checks if short code already exists
func (m *MockRepo) IsShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
	args := m.Called(ctx, shortCode)
	return args.Bool(0), args.Error(1)
}

func TestGenerateRandomCode(t *testing.T) {
	shortCode := generateRandomCode()

	if len(shortCode) != 6 {
		t.Errorf("expected length 6, got %d", len(shortCode))

	}

	for _, char := range shortCode {
		if !strings.ContainsRune(charset, char) {
			t.Errorf("unexpected character %q in code", char)

		}
	}
}

func TestGenerateShortCode(t *testing.T) {
	mockRepo := new(MockRepo)
	service := &URLService{
		repo:    mockRepo,
		baseURL: "https://localhost:8080",
	}

	mockRepo.On("IsShortCodeExists", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)

	shortCode, err := service.GenerateShortCode(context.Background())

	require.NoError(t, err)
	require.NotEmpty(t, shortCode)

}

func TestGenerateShortCode_FailedGeneration(t *testing.T) {
	mockRepo := new(MockRepo)
	service := &URLService{
		repo:    mockRepo,
		baseURL: "https://localhost:8080",
	}

	mockRepo.On("IsShortCodeExists", mock.Anything, mock.AnythingOfType("string")).Return(true, nil)

	shortCode, err := service.GenerateShortCode(context.Background())

	require.Error(t, err)
	require.Empty(t, shortCode)
}

func TestValidateURL(t *testing.T) {

	tests := []struct {
		name          string
		url           string
		err           bool
		expectedError string
	}{
		{
			name:          "valid URL",
			url:           "https://google.com",
			err:           false,
			expectedError: "",
		},
		{
			name:          "empty URL",
			url:           "",
			err:           true,
			expectedError: "URL cannot be empty",
		},
		{
			name:          "no Scheme",
			url:           "google.com",
			err:           true,
			expectedError: "URL must have a scheme (http or https)",
		},
		{
			name:          "invalid domain",
			url:           "https://google",
			err:           true,
			expectedError: "URL must contain a valid domain",
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepo)
			service := &URLService{
				repo:    mockRepo,
				baseURL: "https://localhost:8080",
			}
			err := service.validateURL(tc.url)

			if tc.err != false {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})

	}
}

func TestHealthCheck(t *testing.T) {

	repo := new(MockRepo)
	service := &URLService{
		repo:    repo,
		baseURL: "https://localhost:8080",
	}

	healthRequest := &pb.HealthRequest{}

	response, err := service.HealthCheck(context.Background(), healthRequest)

	require.NoError(t, err)
	require.NotNil(t, response)
	require.True(t, response.Healthy)

}

func TestGetOriginalURL(t *testing.T) {

	repo := new(MockRepo)
	service := &URLService{
		repo:    repo,
		baseURL: "https://localhost:8080",
	}

	request := &pb.GetURLRequest{
		ShortCode: "abc123",
	}

	expectedURL := &models.URL{
		OriginalURL: "https://google.com",
	}

	repo.On("GetByShortCode", mock.Anything, "abc123").Return(expectedURL, nil)

	response, err := service.GetOriginalURL(context.Background(), request)

	require.NoError(t, err)
	require.NotNil(t, response)
	require.Contains(t, response.OriginalUrl, "https://google.com")
	require.True(t, response.Found)

	// invalid shortcode

	request = &pb.GetURLRequest{
		ShortCode: "",
	}

	response, err = service.GetOriginalURL(context.Background(), request)

	require.Error(t, err)
	require.Nil(t, response)
	require.Contains(t, err.Error(), "short_code cannot be empty")

}

func TestCreateShortURL(t *testing.T) {

	repo := new(MockRepo)
	service := &URLService{
		repo:    repo,
		baseURL: "https://localhost:8080/",
		codeGenerator: func(ctx context.Context) (string, error) {
			return "abc123", nil
		},
	}

	request := &pb.CreateURLRequest{
		OriginalUrl: "https://google.com",
		UserId:      "user123",
	}

	repo.On("IsShortCodeExists", mock.Anything, mock.AnythingOfType("string")).Return(false, nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*models.URL")).Return(nil)

	response, err := service.CreateShortURL(context.Background(), request)

	require.NoError(t, err)
	require.NotNil(t, response)
	require.True(t, response.Success)
	require.Equal(t, "abc123", response.ShortCode)
	require.Equal(t, "https://localhost:8080/abc123", response.ShortUrl)

	request = &pb.CreateURLRequest{
		OriginalUrl: "https://",
		UserId:      "user123",
	}

	response, err = service.CreateShortURL(context.Background(), request)

	require.Error(t, err)
	require.Nil(t, response)
	require.Contains(t, err.Error(), "invalid URL:")

	repo = new(MockRepo)
	callCount := 0
	repo.On("IsShortCodeExists", mock.Anything, "abc123").Return(true, nil)

	service = &URLService{
		repo:    repo,
		baseURL: "https://localhost:8080/",
		codeGenerator: func(ctx context.Context) (string, error) {
			for i := 0; i < 10; i++ {
				callCount++
				exists, err := repo.IsShortCodeExists(ctx, "abc123")
				if err != nil {
					return "", err
				}
				if !exists {
					return "abc123", nil
				}
			}
			return "", fmt.Errorf("failed to generate a unique short code after 10 attempts")
		},
	}

	request = &pb.CreateURLRequest{
		OriginalUrl: "https://google.com",
		UserId:      "user123",
	}

	response, err = service.CreateShortURL(context.Background(), request)

	require.Error(t, err)
	require.Nil(t, response)
	require.Contains(t, err.Error(), "failed to generate")
	require.Equal(t, 10, callCount)
}
