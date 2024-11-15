// Package cognito is used as a wrapper for all of our AWS cognito functions.
package s3

import (
	"context"
	"organize-this/infra/logger"
	"sync"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	// Client is a singleton s3 client connection
	client        S3Client
	presignClient S3PresignClient
	once          sync.Once
	err           error
)

type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

type S3PresignClient interface {
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

func S3ClientInit() error {
	var err error
	once.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			logger.Fatalf("Issue creating S3 Client: %v", err)
		}

		newClient := s3.NewFromConfig(cfg)

		client = newClient
		presignClient = s3.NewPresignClient(newClient)
	})

	return err
}

func GetClient() S3Client {
	return client
}

func GetPresignClient() S3PresignClient {
	return presignClient
}
