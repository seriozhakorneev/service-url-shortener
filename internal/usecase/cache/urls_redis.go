package cache

import (
	"errors"
	"fmt"
	"time"

	goRedis "github.com/go-redis/redis"

	"service-url-shortener/pkg/redis"
)

// UrlsCache -.
type UrlsCache struct {
	*redis.Redis
}

// New -.
func New(rd *redis.Redis) *UrlsCache {
	return &UrlsCache{rd}
}

// Set -.
func (c *UrlsCache) Set(key, value string, ttl time.Duration) error {
	err := c.Client.Set(key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("UrlsCache - Set - c.Client.Set: %w", err)
	}
	return nil
}

// Get -.
func (c *UrlsCache) Get(key string) (*string, error) {
	value, err := c.Client.Get(key).Result()
	if err != nil {
		if errors.Is(err, goRedis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return &value, nil
}
