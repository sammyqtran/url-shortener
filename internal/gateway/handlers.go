package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/sammyqtran/url-shortener/internal/queue"
	pb "github.com/sammyqtran/url-shortener/proto"
	"go.uber.org/zap"
)

type GatewayServer struct {
	GrpcClient pb.URLServiceClient
	Publisher  queue.EventPublisher
	Logger     *zap.Logger
}

func (s *GatewayServer) HandleCreateShortURL(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("Incoming request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("client_ip", s.getClientIP(r)),
	)
	defer r.Body.Close()

	var req struct {
		URL string `json:"url"`
	}

	jsonErr := json.NewDecoder(r.Body).Decode(&req)

	if jsonErr != nil {
		s.Logger.Error("Failed to decode JSON request", zap.Error(jsonErr))
		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	if req.URL == "" {
		s.Logger.Warn("Empty URL in request")
		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	request := &pb.CreateURLRequest{
		OriginalUrl: req.URL,
		UserId:      "abc123",
	}

	response, err := s.GrpcClient.CreateShortURL(ctx, request)

	if err != nil {
		s.Logger.Error("gRPC CreateShortURL failed", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Failed to create short URL")
		return
	}

	// Publish URL created event
	if s.Publisher != nil {
		go func() {
			ctx := context.Background()
			err := s.Publisher.PublishURLCreated(ctx, response.ShortCode, req.URL, s.getClientInfo(r))
			if err != nil {
				s.Logger.Error("Failed to publish URL created event", zap.Error(err))
			}
		}()
	}

	resp := map[string]string{"shortcode": response.ShortCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func (s *GatewayServer) HandleGetOriginalURL(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info("Incoming request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("client_ip", s.getClientIP(r)),
	)
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" || shortCode == "create" || shortCode == "healthz" {
		s.Logger.Warn("Invalid shortCode path requested", zap.String("shortCode", shortCode))
		http.NotFound(w, r)
		return
	}

	request := &pb.GetURLRequest{
		ShortCode: shortCode,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	response, err := s.GrpcClient.GetOriginalURL(ctx, request)

	if err != nil {
		s.Logger.Error("gRPC GetOriginalURL failed", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !response.Found {
		http.NotFound(w, r)
		return
	}

	// Publish URL accessed event
	if s.Publisher != nil {
		go func() {
			ctx := context.Background()
			err := s.Publisher.PublishURLAccessed(
				ctx,
				shortCode,
				response.OriginalUrl,
				r.UserAgent(),
				s.getClientIP(r),
				r.Header.Get("Referer"),
			)
			if err != nil {
				s.Logger.Error("Failed to publish URL accessed event", zap.Error(err))
			}
		}()
	}

	http.Redirect(w, r, response.OriginalUrl, http.StatusFound)

}

func (s *GatewayServer) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func respondWithError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (s *GatewayServer) getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}

// getClientInfo extracts client information from request (for analytics)
func (g *GatewayServer) getClientInfo(r *http.Request) string {
	// Extract user info if available (e.g., from JWT token, API key, etc.)
	// For now, return IP address
	return g.getClientIP(r)
}
