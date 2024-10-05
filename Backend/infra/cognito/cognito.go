// Package cognito is used as a wrapper for all of our AWS cognito functions.
package cognito

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

var (
	// Client is a singleton redis client connection
	client *cognitoidentityprovider.Client
	once   sync.Once
	err    error
)

func CognitoClient() error {
	var err error
	once.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		client = cognitoidentityprovider.NewFromConfig(cfg)
	})

	return err
}

func GetClient() *cognitoidentityprovider.Client {
	return client
}
