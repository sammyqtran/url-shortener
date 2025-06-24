// cmd/test-client/main.go
package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/sammyqtran/url-shortener/proto"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewURLServiceClient(conn)
	ctx := context.Background()

	// Test 1: Health Check
	log.Println("=== Testing Health Check ===")
	healthResp, err := client.HealthCheck(ctx, &pb.HealthRequest{})
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	log.Printf("Health status: %v\n", healthResp.Healthy)

	// Test 2: Create Short URL
	log.Println("\n=== Testing Create Short URL ===")
	createReq := &pb.CreateURLRequest{
		OriginalUrl: "https://www.google.com",
		UserId:      "user123",
	}

	createResp, err := client.CreateShortURL(ctx, createReq)
	if err != nil {
		log.Fatalf("Create URL failed: %v", err)
	}

	if !createResp.Success {
		log.Fatalf("Create URL unsuccessful: %s", createResp.Error)
	}

	log.Printf("Short code: %s", createResp.ShortCode)
	log.Printf("Short URL: %s", createResp.ShortUrl)
	log.Printf("Success: %v", createResp.Success)

	// Test 3: Get Original URL
	log.Println("\n=== Testing Get Original URL ===")
	getReq := &pb.GetURLRequest{
		ShortCode: createResp.ShortCode,
	}

	getResp, err := client.GetOriginalURL(ctx, getReq)
	if err != nil {
		log.Fatalf("Get URL failed: %v", err)
	}

	if !getResp.Found {
		log.Fatalf("URL not found: %s", getResp.Error)
	}

	log.Printf("Retrieved URL: %s", getResp.OriginalUrl)
	log.Printf("Found: %v", getResp.Found)

	// Test 4: Access URL multiple times to test click tracking
	log.Println("\n=== Testing Click Tracking ===")
	for i := 0; i < 3; i++ {
		time.Sleep(100 * time.Millisecond) // Small delay for async click counting

		getResp, err := client.GetOriginalURL(ctx, getReq)
		if err != nil {
			log.Fatalf("Get URL failed: %v", err)
		}

		if getResp.Found {
			log.Printf("Access %d - URL: %s", i+2, getResp.OriginalUrl)
		} else {
			log.Printf("Access %d - Error: %s", i+2, getResp.Error)
		}
	}

	// Test 5: Test invalid short code
	log.Println("\n=== Testing Invalid Short Code ===")
	invalidReq := &pb.GetURLRequest{
		ShortCode: "nonexistent",
	}

	_, err = client.GetOriginalURL(ctx, invalidReq)
	if err != nil {
		log.Printf("Expected error for invalid code: %v", err)
	} else {
		// Check if the response indicates not found
		invalidResp, _ := client.GetOriginalURL(ctx, invalidReq)
		if !invalidResp.Found {
			log.Printf("Expected 'not found' for invalid code: %s", invalidResp.Error)
		}
	}

	log.Println("\n=== All tests completed ===")
}
