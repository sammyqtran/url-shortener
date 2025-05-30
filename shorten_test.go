package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {

	// Global setup
	codeToURL = make(map[string]string)
	URLtoCode = make(map[string]string)

	code := m.Run()

	// Global teardown if needed
	os.Exit(code)

}

func TestPingHandler(t *testing.T) {

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	pingHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got status code %d", http.StatusOK, resp.StatusCode)
	}

	var response map[string]string

	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok' got '%s'", response["status"])
	}

}

func TestGenerateShortCode(t *testing.T) {

	length := 6
	shortCode := generateShortCode(length)

	if len(shortCode) != length {
		t.Errorf("expected length %d, got %d", length, len(shortCode))
	}

	for _, c := range shortCode {

		if !strings.ContainsRune(charset, c) {
			t.Errorf("unexpected character %q in code", c)

		}

	}

	code2 := generateShortCode(length)
	if shortCode == code2 {
		t.Logf("warning: two generated codes are equal (could happen but unlikely)")
	}

	shortCode = generateShortCode(0)

	if len(shortCode) != 0 {
		t.Errorf("expected length %d, got %d", 0, len(shortCode))
	}

}

// need to test for these:
// valid url,  empty url, malformed json, existing shortcode,
func TestPostHandler(t *testing.T) {

	type urlRequest struct {
		URL string `json:"url"`
	}
	type urlResponse struct {
		ShortCode string `json:"short_code"`
	}

	tests := []struct {
		name           string
		inputJSON      string
		expectedStatus int
		expectError    bool
		setup          func()
	}{
		{
			name:           "valid URL",
			inputJSON:      `{"url": "https://example.com"}`,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "missing URL field",
			inputJSON:      `{}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "malformed JSON",
			inputJSON:      `{"url":`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "duplicate URL",
			inputJSON:      `{"url": "https://example.com"}`,
			expectedStatus: http.StatusOK,
			expectError:    false,
			setup: func() {
				codeToURL["existingcode"] = "https://example.com"
				URLtoCode["https://example.com"] = "existingcode"
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			req := httptest.NewRequest(http.MethodPost, "/post", strings.NewReader(tc.inputJSON))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			postHandler(rr, req)

			if (rr.Code) != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectError {
				return // no need to continue checks on error cases
			}

			if rr.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("expected content-type application/json got %s", rr.Header().Get("Content-Type"))
			}

			var jsonResponse urlResponse
			err := json.NewDecoder(rr.Body).Decode(&jsonResponse)
			if err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			var reqData urlRequest
			err = json.Unmarshal([]byte(tc.inputJSON), &reqData)
			if err != nil {
				t.Fatalf("failed to unmarshal inputJSON: %v", err)
			}

			mappedURL, exists := codeToURL[jsonResponse.ShortCode]
			if !exists {
				t.Fatalf("shortcode %s not found in codeToURL map", jsonResponse.ShortCode)
			}
			if mappedURL != reqData.URL {
				t.Errorf("expected mapped URL %s, got %s", reqData.URL, mappedURL)
			}

			t.Cleanup(func() { codeToURL = map[string]string{}; URLtoCode = map[string]string{} })
		})

	}

}

func TestGetHandler(t *testing.T) {

	tests := []struct {
		name           string
		shortcode      string
		originalURL    string
		expectedStatus int
		expectError    bool
		expectRedirect bool
	}{
		{
			name:           "valid shortcode",
			shortcode:      "abc123",
			originalURL:    "https://example.com",
			expectedStatus: http.StatusFound,
			expectError:    false,
			expectRedirect: true,
		},
		{
			name:           "empty shortcode",
			shortcode:      "",
			originalURL:    "",
			expectedStatus: http.StatusNotFound,
			expectError:    false,
			expectRedirect: false,
		},
		{
			name:           "invalid shortcode",
			shortcode:      "def456",
			originalURL:    "",
			expectedStatus: http.StatusNotFound,
			expectError:    false,
			expectRedirect: false,
		},
	}

	codeToURL = map[string]string{
		"abc123": "https://example.com",
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/get/"+tc.shortcode, nil)
			w := httptest.NewRecorder()

			getHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectRedirect {
				loc, err := resp.Location()

				if err != nil || loc.String() != tc.originalURL {
					t.Errorf("Expected redirect to %s, got %s", tc.originalURL, loc.String())
				}
			}
			t.Cleanup(func() { codeToURL = map[string]string{}; URLtoCode = map[string]string{} })
		})

	}

}

// codeToURL = map[string]string{
// 	"abc123": "https://example.com",
// }

// URLtoCode = map[string]string{
// 	"https://example.com": "abc123",
// }

// inputJSON := `{"url": "https://example.com"}`

// t.Run("Duplicate URL", func(t *testing.T) {
// 	req := httptest.NewRequest(http.MethodPost, "/post", strings.NewReader(inputJSON))
// 	req.Header.Set("Content-Type", "application/json")
// 	rr := httptest.NewRecorder()

// 	postHandler(rr, req)

// 	if rr.Header().Get("Content-Type") != "application/json" {
// 		t.Fatalf("expected content-type application/json got %s", rr.Header().Get("Content-Type"))
// 	}

// 	if (rr.Code) != http.StatusOK {
// 		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
// 	}
// 	t.Cleanup(func() { codeToURL = map[string]string{}; URLtoCode = map[string]string{} })
// })
