package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	pb "github.com/sammyqtran/url-shortener/proto"
)

type GatewayServer struct {
	GrpcClient pb.URLServiceClient
}

func (s *GatewayServer) HandleCreateShortURL(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var req struct {
		URL string `json:"url"`
	}

	jsonErr := json.NewDecoder(r.Body).Decode(&req)

	if jsonErr != nil || req.URL == "" {
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
		respondWithError(w, http.StatusInternalServerError, "Failed to create short URL")
		return
	}

	resp := map[string]string{"shortcode": response.ShortCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func (s *GatewayServer) HandleGetOriginalURL(w http.ResponseWriter, r *http.Request) {

	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" || shortCode == "create" || shortCode == "healthz" {
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
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !response.Found {
		http.NotFound(w, r)
		return
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
