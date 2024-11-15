package config

import (
	"github.com/spf13/viper"
)

func AWSRegion() string {
	return viper.GetString("AWS_REGION")
}
