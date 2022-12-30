// Package redis implements redis client connection.
package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Redis -.
type Redis struct {
	client *redis.Client
}

func (r *Redis) status() error {
	answer, err := r.client.Ping().Result()
	if err != nil {
		return fmt.Errorf("status - r.client.Ping.Result: %w", err)
	}

	if answer != "PONG" {
		return fmt.Errorf(
			"status - r.client..Ping.Result"+
				"connection is not set, expected: 'PONG', got: %s",
			answer,
		)
	}

	fmt.Println(answer)
	return nil
}

// New -.
func New(port, pass string, db int, opts ...Option) (*Redis, error) {
	rd := &Redis{
		client: redis.NewClient(
			&redis.Options{
				Addr:     "localhost:" + port,
				Password: pass,
				DB:       db,
			},
		),
	}

	// Custom options
	for _, opt := range opts {
		opt(rd)
	}

	err := rd.status()
	if err != nil {
		return nil, fmt.Errorf("redis - NewRedis - rd.Status: %w", err)
	}

	return rd, nil
}

// Close -.
func (r *Redis) Close() error {
	err := r.client.Close()
	if err != nil {
		return fmt.Errorf("redis - Close - r.client.Close: %w", err)
	}

	return nil
}
