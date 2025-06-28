package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/sammyqtran/url-shortener/internal/models"
	"github.com/sammyqtran/url-shortener/internal/repository"
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

	if url, ok := args.Get(0).(*models.URL); ok {
		return url, args.Error(1)
	}
	return nil, args.Error(1)
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

	mockRepo = new(MockRepo)
	service = &URLService{
		repo:    mockRepo,
		baseURL: "https://localhost:8080",
	}

	mockRepo.On("IsShortCodeExists", mock.Anything, mock.AnythingOfType("string")).Return(false, fmt.Errorf("random error"))
	shortCode, err = service.GenerateShortCode(context.Background())

	require.Error(t, err)
	require.Empty(t, shortCode)
	require.Contains(t, err.Error(), "error checking short code existence:")
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
		{
			name:          "invalid URL format",
			url:           "http://%41:8080/",
			err:           true,
			expectedError: "invalid URL format:",
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
	tests := []struct {
		name          string
		mockReturn    interface{}
		mockError     error
		request       *pb.GetURLRequest
		expectError   bool
		expectNilResp bool
		checkResponse func(t *testing.T, resp *pb.GetURLResponse, err error)
		hitCache      bool
	}{
		{
			name: "success",
			mockReturn: &models.URL{
				OriginalURL: "https://google.com",
			},
			mockError: nil,
			request: &pb.GetURLRequest{
				ShortCode: "abc123",
			},
			expectError:   false,
			expectNilResp: false,
			checkResponse: func(t *testing.T, resp *pb.GetURLResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Contains(t, resp.OriginalUrl, "https://google.com")
				require.True(t, resp.Found)
			},
		},
		{
			name:          "invalid shortcode empty",
			mockReturn:    nil,
			mockError:     nil,
			request:       &pb.GetURLRequest{ShortCode: ""},
			expectError:   true,
			expectNilResp: true,
			checkResponse: func(t *testing.T, resp *pb.GetURLResponse, err error) {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Contains(t, err.Error(), "short_code cannot be empty")
			},
		},
		{
			name:       "repo url not found",
			mockReturn: nil,
			mockError:  repository.ErrURLNotFound,
			request: &pb.GetURLRequest{
				ShortCode: "abc123",
			},
			expectError:   false,
			expectNilResp: false,
			checkResponse: func(t *testing.T, resp *pb.GetURLResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.False(t, resp.Found)
				require.Equal(t, "URL not found", resp.Error)
			},
		},
		{
			name:       "repo other failure",
			mockReturn: nil,
			mockError:  fmt.Errorf("random failure"),
			request: &pb.GetURLRequest{
				ShortCode: "abc123",
			},
			expectError:   false,
			expectNilResp: false,
			checkResponse: func(t *testing.T, resp *pb.GetURLResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.False(t, resp.Found)
				require.Contains(t, resp.Error, "failed to retrieve URL:")
			},
		},
		{
			name: "expired URL",
			mockReturn: &models.URL{
				OriginalURL: "https://google.com",
				ExpiresAt:   ptrTime(time.Now().Add(-time.Hour)),
			},
			mockError: nil,
			request: &pb.GetURLRequest{
				ShortCode: "abc123",
			},
			expectError:   false,
			expectNilResp: false,
			checkResponse: func(t *testing.T, resp *pb.GetURLResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.False(t, resp.Found)
				require.Contains(t, resp.Error, "URL has expired")
			},
		},
		{
			name: "db fetch populates cache",
			mockReturn: &models.URL{
				OriginalURL: "https://example.com",
				ExpiresAt:   nil, // or future time, so not expired
			},
			mockError: nil,
			request: &pb.GetURLRequest{
				ShortCode: "abc123",
			},
			expectError:   false,
			expectNilResp: false,
			checkResponse: func(t *testing.T, resp *pb.GetURLResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.True(t, resp.Found)
				require.Equal(t, "https://example.com", resp.OriginalUrl)
			},
			hitCache: false, // force cache miss so it falls back to repo
		},
		{
			name: "expired in cache",
			mockReturn: &models.URL{
				OriginalURL: "https://google.com",
				ExpiresAt:   ptrTime(time.Now().Add(-time.Hour)),
			},
			mockError: nil,
			request: &pb.GetURLRequest{
				ShortCode: "abc123",
			},
			expectError:   false,
			expectNilResp: false,
			checkResponse: func(t *testing.T, resp *pb.GetURLResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.False(t, resp.Found)
				require.Contains(t, resp.Error, "URL has expired")
			},
			hitCache: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mockClient := redismock.NewClientMock()
			repo := new(MockRepo)
			service := &URLService{
				repo:    repo,
				baseURL: "https://localhost:8080",
				cache:   db,
			}

			if tt.mockReturn != nil || tt.mockError != nil {
				repo.On("GetByShortCode", mock.Anything, mock.Anything).Return(tt.mockReturn, tt.mockError)
			}

			if tt.name == "success" {
				data, _ := json.Marshal(tt.mockReturn)
				mockClient.ExpectGet("url:" + tt.request.ShortCode).SetVal(string(data))
			} else {
				if tt.hitCache {
					data, _ := json.Marshal(tt.mockReturn)
					mockClient.ExpectGet("url:" + tt.request.ShortCode).SetVal(string(data))
					mockClient.ExpectDel("url:" + tt.request.ShortCode).SetVal(1)
				} else {
					mockClient.ExpectGet("url:" + tt.request.ShortCode).RedisNil()
				}
			}

			resp, err := service.GetOriginalURL(context.Background(), tt.request)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if tt.expectNilResp {
				require.Nil(t, resp)
			} else {
				require.NotNil(t, resp)
			}

			tt.checkResponse(t, resp, err)
		})
	}
}

func TestCreateShortURL(t *testing.T) {
	tests := []struct {
		name          string
		codeGenerator func(ctx context.Context) (string, error)
		mockSetup     func(m *MockRepo, mockRedis redismock.ClientMock)
		request       *pb.CreateURLRequest
		expectError   bool
		checkResponse func(t *testing.T, resp *pb.CreateURLResponse, err error)
	}{
		{
			name: "success",
			codeGenerator: func(ctx context.Context) (string, error) {
				return "abc123", nil
			},
			mockSetup: func(m *MockRepo, mockRedis redismock.ClientMock) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*models.URL")).Return(nil)
				mockRedis.ExpectSet("abc123", mock.Anything, 10*time.Minute).SetVal("OK")

			},
			request: &pb.CreateURLRequest{
				OriginalUrl: "https://google.com",
				UserId:      "user123",
			},
			expectError: false,
			checkResponse: func(t *testing.T, resp *pb.CreateURLResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.True(t, resp.Success)
				require.Equal(t, "abc123", resp.ShortCode)
				require.Equal(t, "https://localhost:8080/abc123", resp.ShortUrl)
			},
		},
		{
			name:      "invalid url",
			mockSetup: func(m *MockRepo, mockRedis redismock.ClientMock) {},
			codeGenerator: func(ctx context.Context) (string, error) {
				return "abc123", nil
			},
			request: &pb.CreateURLRequest{
				OriginalUrl: "https://",
				UserId:      "user123",
			},
			expectError: true,
			checkResponse: func(t *testing.T, resp *pb.CreateURLResponse, err error) {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Contains(t, err.Error(), "invalid URL:")
			},
		},
		{
			name: "failed code generation",
			codeGenerator: func(ctx context.Context) (string, error) {
				return "", fmt.Errorf("failed to generate a unique short code after 10 attempts")
			},
			mockSetup: func(m *MockRepo, mockRedis redismock.ClientMock) {
			},
			request: &pb.CreateURLRequest{
				OriginalUrl: "https://google.com",
				UserId:      "user123",
			},
			expectError: true,
			checkResponse: func(t *testing.T, resp *pb.CreateURLResponse, err error) {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Contains(t, err.Error(), "failed to generate a unique short code")
			},
		},
		{
			name: "failed save to db",
			codeGenerator: func(ctx context.Context) (string, error) {
				return "abc123", nil
			},
			mockSetup: func(m *MockRepo, mockRedis redismock.ClientMock) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*models.URL")).Return(fmt.Errorf("random failure"))
			},
			request: &pb.CreateURLRequest{
				OriginalUrl: "https://google.com",
				UserId:      "user123",
			},
			expectError: true,
			checkResponse: func(t *testing.T, resp *pb.CreateURLResponse, err error) {
				require.Error(t, err)
				require.Nil(t, resp)
				require.Contains(t, err.Error(), "failed to create URL:")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRepo)

			cache, mockRedis := redismock.NewClientMock()
			tt.mockSetup(repo, mockRedis)

			service := &URLService{
				repo:          repo,
				baseURL:       "https://localhost:8080/",
				codeGenerator: tt.codeGenerator,
				cache:         cache,
			}

			resp, err := service.CreateShortURL(context.Background(), tt.request)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			tt.checkResponse(t, resp, err)
			repo.AssertExpectations(t)
		})
	}
}

func TestGetFromCache_CacheHit(t *testing.T) {
	// create redis mock client and mock controller
	db, mock := redismock.NewClientMock()
	repo := new(MockRepo)
	service := &URLService{
		repo:    repo,
		cache:   db,
		baseURL: "https://localhost:8080",
	}

	urlModel := &models.URL{
		OriginalURL: "https://google.com",
		ShortCode:   "abc123",
	}

	data, err := json.Marshal(urlModel)
	require.NoError(t, err)

	// mock redis GET command expectation for the key "url:abc123"
	mock.ExpectGet("url:abc123").SetVal(string(data))

	result, err := service.getFromCache(context.Background(), "abc123")
	require.NoError(t, err)
	require.Equal(t, urlModel.OriginalURL, result.OriginalURL)

	require.NoError(t, mock.ExpectationsWereMet())

}

func TestGetFromCache_CacheMiss(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := &URLService{
		repo:  new(MockRepo),
		cache: db,
	}

	mock.ExpectGet("url:abc123").RedisNil()

	result, err := service.getFromCache(context.Background(), "abc123")
	require.Error(t, err)
	require.Nil(t, result)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFromCache_CachErr(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := &URLService{
		repo:  new(MockRepo),
		cache: db,
	}

	mock.ExpectGet("url:abc123").SetErr(fmt.Errorf("redis error"))

	result, err := service.getFromCache(context.Background(), "abc123")
	require.Error(t, err)
	require.Nil(t, result)

}

func TestSetCacheFromModel(t *testing.T) {
	db, mock := redismock.NewClientMock()
	service := &URLService{
		repo:  new(MockRepo),
		cache: db,
	}

	urlModel := &models.URL{
		OriginalURL: "https://google.com",
		ShortCode:   "abc123",
	}

	data, err := json.Marshal(urlModel)
	require.NoError(t, err)

	mock.ExpectSet("url:abc123", data, 10*time.Minute).SetVal("OK")

	err = service.setCacheFromModel(context.Background(), "abc123", urlModel)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
