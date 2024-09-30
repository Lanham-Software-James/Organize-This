package cache

import (
	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	err    error
)

// ClientConnection create redis connection
func ClientConnection(redisConnectionString string) error {
	var client = Client

	opt, err := redis.ParseURL(redisConnectionString)
	if err != nil {
		panic(err)
	}

	client = redis.NewClient(opt)

	Client = client

	return nil
}

// GetClient Redis connection
func GetClient() *redis.Client {
	return Client
}
