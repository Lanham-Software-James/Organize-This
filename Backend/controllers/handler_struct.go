// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"organize-this/helpers"
	"organize-this/infra/cognito"
	"organize-this/repository"
)

// Handler is used to allow us to pass our database to the controllers enabling us to mock during unit testing.
type Handler struct {
	Repository    *repository.Repository
	CognitoClient cognito.CognitoClient
	TokenHelper   helpers.TokenHelper
}
