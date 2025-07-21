package main

import (
	"net/http"
	"time"
	"willowsuite-vault/config"
	"willowsuite-vault/infra/cache"
	"willowsuite-vault/infra/cognito"
	"willowsuite-vault/infra/database"
	"willowsuite-vault/infra/logger"
	"willowsuite-vault/infra/s3"
	"willowsuite-vault/migrations"
	"willowsuite-vault/routers"

	"github.com/spf13/viper"
)

func main() {

	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	masterDSN, replicaDSN := config.DbConfiguration()
	if err := database.DbConnection(masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}
	//later separate migration
	migrations.Migrate()

	redisConnectionString := config.RedisConfiguration()
	if err := cache.ClientConnection(redisConnectionString); err != nil {
		logger.Fatalf("redis ClientConnection error: %s", err)
	}

	if err := cognito.CognitoClientInit(); err != nil {
		logger.Fatalf("Cognito Connection error: %s", err)
	}

	if err := s3.S3ClientInit(); err != nil {
		logger.Fatalf("S3 Connection error: %s", err)
	}

	router := routers.SetupRoute()
	logger.Fatalf("%v", http.ListenAndServe(config.ServerConfig(), router))

}
