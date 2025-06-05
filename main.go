package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	//setup db

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		log.Fatal("REDIS_ADDR environment variable is not set")
	}
	initRedis(addr)

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/get/", getHandler)
	http.HandleFunc("/post", postHandler)
	http.ListenAndServe(":8080", nil)
}
