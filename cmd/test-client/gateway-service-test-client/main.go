//cmd/test-client/gateway-service-test-client/main.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	baseURL := "http://localhost:8080"

	fmt.Println("[*] Starting health check test...")
	if err := testHealth(baseURL + "/healthz"); err != nil {
		fail(err)
	}
	fmt.Println("[+] Health check passed")

	fmt.Println("[*] Starting create short URL test...")
	shortCode, err := testCreate(baseURL + "/create")
	if err != nil {
		fail(err)
	}
	fmt.Printf("[+] Create short URL passed, shortcode: %s\n", shortCode)

	fmt.Println("[*] Starting get original URL test...")
	if err := testGet(baseURL + "/" + shortCode); err != nil {
		fail(err)
	}
	fmt.Println("[+] Get original URL test passed")

	fmt.Println("[+] All tests passed!")
}

func testHealth(url string) error {
	fmt.Printf("-> Sending GET request to %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("<- Received status code: %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("health check returned status %d, expected 200; body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func testCreate(url string) (string, error) {
	payload := map[string]string{"url": "https://example.com"}
	body, _ := json.Marshal(payload)

	fmt.Printf("-> Sending POST request to %s with body: %s\n", url, string(body))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("<- Received status code: %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("create returned status %d, expected 200; body: %s", resp.StatusCode, string(body))
	}

	var data struct {
		Shortcode string `json:"shortcode"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode create response: %w", err)
	}

	fmt.Printf("<- Received shortcode: %s\n", data.Shortcode)
	if data.Shortcode == "" {
		return "", fmt.Errorf("empty shortcode in create response")
	}

	return data.Shortcode, nil
}

func testGet(url string) error {
	fmt.Printf("-> Sending GET request to %s\n", url)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Prevent following redirects
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("get request failed: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("<- Received status code: %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("get returned status %d, expected 302; body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func fail(err error) {
	fmt.Println("[!] Test failed:", err)
	os.Exit(1)
}
