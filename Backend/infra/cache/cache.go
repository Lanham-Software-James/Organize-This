// Package cache is the homebase for our redis connection
package cache

import (
	"organize-this/infra/logger"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	// Client is a singleton redis client connection
	Client *redis.Client
	once   sync.Once
	err    error
)

// ClientConnection create redis connection
func ClientConnection(redisConnectionString string) error {
	var client = Client
	var err error
	once.Do(func() {
		opt, err := redis.ParseURL(redisConnectionString)
		if err != nil {
			logger.Fatalf("error connecting to redis: %v", err)
		}

		client = redis.NewClient(opt)
		Client = client
	})

	return err
}

// GetClient Redis connection
func GetClient() *redis.Client {
	return Client
}
