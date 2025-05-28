package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
)

// map holding url mapping
var codeToURL = make(map[string]string)
var URLtoCode = make(map[string]string)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func pingHandler(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func getHandler(w http.ResponseWriter, r *http.Request) {

	shortURL := strings.TrimPrefix(r.URL.Path, "/get/")
	originalURL, exists := codeToURL[shortURL]

	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)

}

func postHandler(w http.ResponseWriter, r *http.Request) {

	var req struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil || req.URL == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
		return
	}

	originalURL := req.URL

	if existingCode, found := URLtoCode[originalURL]; found {
		resp := map[string]string{"short_code": existingCode}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	var shortCode string

	for {
		shortCode = generateShortCode(6)
		if _, found := codeToURL[shortCode]; !found {
			break
		}

	}

	codeToURL[shortCode] = originalURL
	URLtoCode[originalURL] = shortCode

	resp := map[string]string{"short_code": shortCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}

func generateShortCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)

}

func main() {

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/get/", getHandler)
	http.HandleFunc("/post", postHandler)
	http.ListenAndServe(":8080", nil)
}
