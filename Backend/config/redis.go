package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	Address  string
	Password string
}

func RedisConfiguration() string {
	redisAddress := viper.GetString("REDIS_HOST")
	redisUser := viper.GetString("REDIS_USER")
	redisPassword := viper.GetString("REDIS_PASSWORD")

	redisConnectionString := fmt.Sprintf(
		"redis://%s:%s@%s/0",
		redisUser, redisPassword, redisAddress,
	)

	return redisConnectionString
}
