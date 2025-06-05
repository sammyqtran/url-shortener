package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func initRedis(addr string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", addr, err)
	}
}

func saveMapping(code, url string) error {

	shortCodeErr := rdb.Set(ctx, "shortcode:"+code, url, 0).Err()

	if shortCodeErr != nil {
		log.Printf("ERROR: Redis SET failed for key %s : %v", code, shortCodeErr)
		return shortCodeErr
	}

	urlError := rdb.Set(ctx, "url:"+url, code, 0).Err()

	if urlError != nil {
		log.Printf("ERROR: Redis SET failed for key %s: %v", url, urlError)

	}

	return nil
}

func getURL(code string) (string, bool, error) {

	url, err := rdb.Get(ctx, "shortcode:"+code).Result()

	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		log.Printf("ERROR: Redis GET failed for key %s : %v", code, err)
		return "", false, err
	}

	return url, true, nil
}

func getShortCode(url string) (string, bool, error) {

	shortCode, err := rdb.Get(ctx, "url:"+url).Result()

	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		log.Printf("ERROR: Redis GET failed for key %s : %v", url, err)
		return "", false, nil
	}

	return shortCode, true, err

}
