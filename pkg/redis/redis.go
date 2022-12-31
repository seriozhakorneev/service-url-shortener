// Package redis implements redis Client connection.
package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

// Redis -.
type Redis struct {
	Client *redis.Client
}

// New -.
func New(address, pass string, db int, opts ...Option) (*Redis, error) {
	rd := &Redis{
		Client: redis.NewClient(
			&redis.Options{
				Addr:     address,
				Password: pass,
				DB:       db,
			},
		),
	}

	// Custom options
	for _, opt := range opts {
		opt(rd)
	}

	err := rd.ping()
	if err != nil {
		return nil, fmt.Errorf("redis - NewRedis - rd.ping: %w", err)
	}

	return rd, nil
}

func (r *Redis) ping() error {
	answer, err := r.Client.Ping().Result()
	if err != nil {
		return fmt.Errorf("ping - r.Client.Ping.Result: %w", err)
	}

	if answer != "PONG" {
		return fmt.Errorf(
			"ping - r.Client..Ping.Result"+
				"connection is not set, expected: 'PONG', got: %s",
			answer,
		)
	}

	return nil
}

// Close -.
func (r *Redis) Close() error {
	err := r.Client.Close()
	if err != nil {
		return fmt.Errorf("redis - Close - r.Client.Close: %w", err)
	}

	return nil
}
