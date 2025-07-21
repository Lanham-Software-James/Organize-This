// Package controllers provides all the various controllers for requests coming to the WillowSuite Vault API.
package controllers

import (
	"willowsuite-vault/helpers"
	"willowsuite-vault/infra/cognito"
	"willowsuite-vault/infra/s3"
	"willowsuite-vault/repository"
)

// Handler is used to allow us to pass our database to the controllers enabling us to mock during unit testing.
type Handler struct {
	Repository      *repository.Repository
	CognitoClient   cognito.CognitoClient
	S3Client        s3.S3Client
	S3PresignClient s3.S3PresignClient
	TokenHelper     helpers.TokenHelper
}
