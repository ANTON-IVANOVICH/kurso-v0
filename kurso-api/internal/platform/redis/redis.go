// Package redis provides a go-redis client (a driven adapter).
package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// New creates a Redis client from a redis:// connection URL. The connection is
// lazy; use Client.Ping to verify reachability.
func New(url string) (*redis.Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	return redis.NewClient(opt), nil
}
