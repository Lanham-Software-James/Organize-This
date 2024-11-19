package config

import (
	"github.com/spf13/viper"
)

func S3BucketName() string {
	return viper.GetString("AWS_S3_BUCKET_NAME")
}
