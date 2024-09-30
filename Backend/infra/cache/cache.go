// Package cache is the homebase for our redis connection
package cache

import (
	"organize-this/infra/logger"

	"github.com/redis/go-redis/v9"
)

var (
	// Client is a singleton redis client connection
	Client *redis.Client
	err    error
)

// ClientConnection create redis connection
func ClientConnection(redisConnectionString string) error {
	var client = Client

	opt, err := redis.ParseURL(redisConnectionString)
	if err != nil {
		logger.Fatalf("error connecting to redis: %v", err)
		return err
	}

	client = redis.NewClient(opt)
	Client = client

	return nil
}

// GetClient Redis connection
func GetClient() *redis.Client {
	return Client
}
