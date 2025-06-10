package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient() *RedisClient {

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     os.Getenv("REDIS_PASSWORD"), // empty if no password
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	ctx := context.Background()

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		// Don't panic - allow graceful degradation
	}

	return &RedisClient{
		client: rdb,
		ctx:    ctx,
	}

}

func (r *RedisClient) saveMapping(code, url string) error {

	shortCodeErr := r.client.Set(r.ctx, "shortcode:"+code, url, 0).Err()

	if shortCodeErr != nil {
		log.Printf("ERROR: Redis SET failed for key %s : %v", code, shortCodeErr)
		return shortCodeErr
	}

	urlError := r.client.Set(r.ctx, "url:"+url, code, 0).Err()

	if urlError != nil {
		log.Printf("ERROR: Redis SET failed for key %s: %v", url, urlError)

	}

	return nil
}

func (r *RedisClient) getURL(code string) (string, bool, error) {

	url, err := r.client.Get(r.ctx, "shortcode:"+code).Result()

	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		log.Printf("ERROR: Redis GET failed for key %s : %v", code, err)
		return "", false, err
	}

	return url, true, nil
}

func (r *RedisClient) getShortCode(url string) (string, bool, error) {

	shortCode, err := r.client.Get(r.ctx, "url:"+url).Result()

	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		log.Printf("ERROR: Redis GET failed for key %s : %v", url, err)
		return "", false, nil
	}

	return shortCode, true, err

}

func (r *RedisClient) incrementClicks(code string) error {
	key := fmt.Sprintf("clicks:%s", code)
	return r.client.Incr(r.ctx, key).Err()
}

func (r *RedisClient) getClicks(code string) (int64, error) {

	key := fmt.Sprintf("clicks:%s", code)
	clicks, err := r.client.Get(r.ctx, key).Result()

	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	count, err := strconv.ParseInt(clicks, 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *RedisClient) HealthCheck() error {
	return r.client.Ping(r.ctx).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
