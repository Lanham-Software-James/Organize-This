package config

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/spf13/viper"
)

func CognitoClientID() string {
	return viper.GetString("AWS_CLIENT_ID")
}

func CognitoUserPoolID() string {
	return viper.GetString("AWS_USER_POOL_ID")
}

func CognitoRegion() string {
	return viper.GetString("AWS_REGION")
}

func cognitoClientSecret() string {
	return viper.GetString("AWS_CLIENT_SECRET")
}

func CognitoSecretHash(userName string) string {
	mac := hmac.New(sha256.New, []byte(cognitoClientSecret()))
	mac.Write([]byte(userName + CognitoClientID()))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
