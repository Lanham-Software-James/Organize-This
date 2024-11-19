package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Port                 string
	Secret               string
	LimitCountPerRequest int64
}

func ServerConfig() string {

	appServer := fmt.Sprintf("%s:%s", viper.GetString("SERVER_HOST"), viper.GetString("SERVER_PORT"))
	log.Print("Server Running at :", appServer)

	return appServer
}

func FrontEndURL() string {
	return viper.GetString("FRONT_END_URL")
}

func EncryptionSecert() string {
	return viper.GetString("ENCRYPTION_SECERT")
}
