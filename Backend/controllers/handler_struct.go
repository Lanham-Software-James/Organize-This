// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"organize-this/repository"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// Handler is used to allow us to pass our database to the controllers enabling us to mock during unit testing.
type Handler struct {
	Repository    *repository.Repository
	CognitoClient *cognitoidentityprovider.Client
}
