// Package cognito is used as a wrapper for all of our AWS cognito functions.
package cognito

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	// Client is a singleton redis client connection
	client *cognitoidentityprovider.Client
	err    error
)

func CognitoClient() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
		return err
	}

	client = cognitoidentityprovider.NewFromConfig(cfg)
	return nil
}

func GetClient() *cognitoidentityprovider.Client {
	return client
}
