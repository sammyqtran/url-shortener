package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func pingHandler(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func getHandler(w http.ResponseWriter, r *http.Request) {

	// currently getting our short code from the map
	// first part is trimming the string after our /get/
	// lets check our rdb to see if its there

	shortURL := strings.TrimPrefix(r.URL.Path, "/get/")
	originalURL, exists, err := getURL(shortURL)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)

}

func postHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var req struct {
		URL string `json:"url"`
	}

	jsonErr := json.NewDecoder(r.Body).Decode(&req)

	if jsonErr != nil || req.URL == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
		return
	}

	originalURL := req.URL

	existingCode, found, err := getShortCode(originalURL)

	if err != nil {
		log.Printf("Redis error: %v", err)
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
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
		_, found, err := getURL(shortCode)
		if err != nil {
			log.Printf("Redis failure during shortcode check: %v", err)
			http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
			return
		}
		if !found {
			break
		}

	}

	if shortCode == "" {
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	dbErr := saveMapping(shortCode, originalURL)

	if dbErr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to save mapping"})
		return
	}

	resp := map[string]string{"shortcode": shortCode}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

var shortCodeGenerator = func(n int) string {
	return generateShortCode(n)
}

func generateShortCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)

}
