package gateway

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pb "github.com/sammyqtran/url-shortener/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// mock grpc client
type MockURLServiceClient struct {
	pb.URLServiceClient
	mock.Mock
}

func (m *MockURLServiceClient) CreateShortURL(ctx context.Context,
	in *pb.CreateURLRequest, opts ...grpc.CallOption) (*pb.CreateURLResponse, error) {

	args := m.Called(ctx, in, opts)
	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}
	return resp.(*pb.CreateURLResponse), args.Error(1)
}

func (m *MockURLServiceClient) GetOriginalURL(ctx context.Context,
	in *pb.GetURLRequest, opts ...grpc.CallOption) (*pb.GetURLResponse, error) {
	args := m.Called(ctx, in, opts)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}
	return resp.(*pb.GetURLResponse), args.Error(1)

}

func (m *MockURLServiceClient) HealthCheck(ctx context.Context,
	in *pb.HealthRequest, opts ...grpc.CallOption) (*pb.HealthResponse, error) {

	args := m.Called(ctx, in, opts)

	resp := args.Get(0)
	if resp == nil {
		return nil, args.Error(1)
	}
	return resp.(*pb.HealthResponse), args.Error(1)
}

func TestHandlleHealthCheck(t *testing.T) {

	mockClient := new(MockURLServiceClient)
	server := &GatewayServer{
		GrpcClient: mockClient,
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	server.HandleHealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

}

func TestHandleGetOriginalURL(t *testing.T) {

	type testCase struct {
		name           string
		mockResponse   *pb.GetURLResponse
		mockError      error
		shortCode      string
		expectGrpcCall bool
		expectedCode   int
		expectError    bool
		expectedError  string
	}

	tests := []testCase{
		{
			name:      "Retrieve Valid URL",
			shortCode: "abc123",
			mockResponse: &pb.GetURLResponse{
				OriginalUrl: "https://google.com",
				Found:       true,
			},
			mockError:      nil,
			expectGrpcCall: true,
			expectedCode:   http.StatusFound,
		},
		{
			name:           "Missing Short code/ Invalid URL",
			shortCode:      "",
			expectGrpcCall: false,
			expectedCode:   http.StatusNotFound,
		},
		{
			name:           "Failed gRPC call",
			shortCode:      "abc123",
			expectGrpcCall: true,
			mockResponse:   nil,
			mockError:      errors.New("grpc failure"),
			expectError:    true,
			expectedError:  "internal server error",
			expectedCode:   http.StatusInternalServerError,
		},
		{
			name:           "URL not found",
			shortCode:      "abc123",
			expectGrpcCall: true,
			mockResponse: &pb.GetURLResponse{
				Found: false,
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {

			mockClient := new(MockURLServiceClient)
			server := &GatewayServer{
				GrpcClient: mockClient,
			}

			if tc.expectGrpcCall {
				mockClient.On("GetOriginalURL", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockResponse, tc.mockError)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/"+tc.shortCode, nil)
			server.HandleGetOriginalURL(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("expected status %d, got %d", tc.expectedCode, w.Code)
			}

			if tc.expectError && !strings.Contains(w.Body.String(), tc.expectedError) {
				t.Errorf("expected error message to contain %q, got %q", tc.expectedError, w.Body.String())

			}

			if tc.expectGrpcCall {
				mockClient.AssertCalled(t, "GetOriginalURL", mock.Anything, mock.Anything, mock.Anything)
			} else {
				mockClient.AssertNotCalled(t, "GetOriginalURL", mock.Anything, mock.Anything, mock.Anything)
			}
		})

	}

}
func TestHandleCreateShortURL(t *testing.T) {

	type testCase struct {
		name           string
		inputBody      string
		mockResponse   *pb.CreateURLResponse
		mockError      error
		expectedCode   int
		expectedBody   string
		expectGrpcCall bool
	}

	testCases := []testCase{
		{
			name:      "valid test",
			inputBody: `{"url": "https://example.com"}`,
			mockResponse: &pb.CreateURLResponse{
				ShortCode: "abc123",
				Success:   true,
				ShortUrl:  "https://localhost:8080/abc123",
			},
			mockError:      nil,
			expectedCode:   http.StatusOK,
			expectedBody:   `{"shortcode":"abc123"}`,
			expectGrpcCall: true,
		},
		{
			name:           "json parsing issue",
			inputBody:      `{"url:`,
			mockResponse:   nil,
			expectedBody:   `{"error":"bad request"}`,
			expectGrpcCall: false,
			expectedCode:   http.StatusBadRequest,
		},
		{
			name:           "gRPC call fails",
			inputBody:      `{"url": "https://example.com"}`,
			mockResponse:   nil,
			mockError:      errors.New("grpc failure"),
			expectedCode:   http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to create short URL"}`,
			expectGrpcCall: true,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(MockURLServiceClient)
			server := &GatewayServer{
				GrpcClient: mockClient,
			}
			if tc.expectGrpcCall {
				mockClient.
					On("CreateShortURL", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.mockResponse, tc.mockError)
			}

			req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(tc.inputBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.HandleCreateShortURL(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("expected status %d, got %d", tc.expectedCode, w.Code)
			}

			if strings.TrimSpace(w.Body.String()) != tc.expectedBody {
				t.Errorf("expected body %s, got %s", tc.expectedBody, w.Body.String())
			}
		})

	}

}
