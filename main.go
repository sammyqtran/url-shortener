package main

import (
	"log"
	"net/http"
)

func main() {

	// Dependency Injection
	// Singleton of dependency (redis) -> injected where needed (http handlers)
	// avoids global state
	redisClient := NewRedisClient()
	if redisClient == nil {
		log.Fatal("Could not initialize Redis client")
	}
	defer redisClient.client.Close()

	handler := &Handler{Redis: redisClient}

	http.HandleFunc("/ping", handler.pingHandler)
	http.HandleFunc("/get/", handler.getHandler)
	http.HandleFunc("/post", handler.postHandler)
	http.HandleFunc("/analytics/", handler.AnalyticsHandler)
	http.ListenAndServe(":8080", nil)
}
