package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	Redis *RedisClient
}

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

// h *Handler is called a pointer reciever for a method on a struct
func (h *Handler) pingHandler(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) {

	// currently getting our short code from the map
	// first part is trimming the string after our /get/
	// lets check our rdb to see if its there

	shortURL := strings.TrimPrefix(r.URL.Path, "/get/")
	originalURL, exists, err := h.Redis.getURL(shortURL)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if !exists {
		http.NotFound(w, r)
		return
	}

	err = h.Redis.incrementClicks(shortURL)

	if err != nil {
		log.Printf("Redis error: %v", err)
	}

	http.Redirect(w, r, originalURL, http.StatusFound)

}

func (h *Handler) postHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var req struct {
		URL string `json:"url"`
	}

	jsonErr := json.NewDecoder(r.Body).Decode(&req)

	if jsonErr != nil || req.URL == "" {
		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	originalURL := req.URL

	existingCode, found, err := h.Redis.getShortCode(originalURL)

	if err != nil {
		log.Printf("Redis error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if found {

		resp := map[string]string{"shortcode": existingCode}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	var shortCode string

	for index := 0; index < 10; index++ {

		shortCode = shortCodeGenerator(6)
		_, found, err := h.Redis.getURL(shortCode)
		if err != nil {
			log.Printf("Redis failure during shortcode check: %v", err)
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		if !found {
			break
		}

	}

	if shortCode == "" {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	dbErr := h.Redis.saveMapping(shortCode, originalURL)

	if dbErr != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to save mapping")
		return
	}

	resp := map[string]string{"shortcode": shortCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

// Analytics Handlers

func (h *Handler) AnalyticsHandler(w http.ResponseWriter, r *http.Request) {

	shortcode := strings.TrimPrefix(r.URL.Path, "/analytics/")

	if shortcode == "" {
		respondWithError(w, http.StatusBadRequest, "shortcode required")
		return
	}

	clicks, err := h.Redis.getClicks(shortcode)

	if err != nil {
		log.Printf("Redis failure during gathering analytics: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := map[string]interface{}{
		"shortcode": shortcode,
		"clicks":    clicks,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

//Helpers

// http error response formatter
func respondWithError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// short code generation
var shortCodeGenerator = func(n int) string {
	return generateShortCode(n)
}

func generateShortCode(length int) string {
	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}

	return string(b)

}
