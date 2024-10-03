package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"organize-this/config"
	"organize-this/helpers"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/go-chi/chi/v5"
)

// SignUp signs up a user with Amazon Cognito.
func (handler Handler) SignUp(w http.ResponseWriter, request *http.Request) {
	byteData, err := io.ReadAll(request.Body)
	if err != nil {
		logAndRespond(w, "Error parsing request", err)
		return
	}

	var parsedData map[string]string
	if err = json.Unmarshal(byteData, &parsedData); err != nil {
		logAndRespond(w, "Error parsing json", err)
		return
	}

	userEmail, password, firstName, lastName, birthday := parsedData["userEmail"], parsedData["password"], parsedData["firstName"], parsedData["lastName"], parsedData["birthday"]
	if userEmail == "" {
		logAndRespond(w, "Missing user name", nil)
		return
	}

	if password == "" {
		logAndRespond(w, "Missing password", nil)
		return
	}

	if firstName == "" {
		logAndRespond(w, "Missing first name", nil)
		return
	}

	if lastName == "" {
		logAndRespond(w, "Missing last name", nil)
		return
	}

	if birthday == "" {
		logAndRespond(w, "Missing birthday", nil)
		return
	}

	output, err := handler.CognitoClient.SignUp(context.TODO(), &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(config.CognitoClientID()),
		Password: aws.String(password),
		Username: aws.String(userEmail),
		UserAttributes: []types.AttributeType{
			{Name: aws.String("given_name"), Value: aws.String(firstName)},
			{Name: aws.String("family_name"), Value: aws.String(lastName)},
			{Name: aws.String("birthdate"), Value: aws.String(birthday)},
		},
		SecretHash: aws.String(config.CognitoSecretHash(userEmail)),
	})
	if err != nil {
		var invalidPassword *types.InvalidPasswordException
		if errors.As(err, &invalidPassword) {
			logAndRespond(w, *invalidPassword.Message, err)
		} else {
			logAndRespond(w, "Couldn't sign up user", err)
		}

		return
	}

	helpers.SuccessResponse(w, &output)
}

// ConfirmSignUp signs up a user with Amazon Cognito.
func (handler Handler) ConfirmSignUp(w http.ResponseWriter, request *http.Request) {
	userEmail := chi.URLParam(request, "email")
	if userEmail == "" {
		logAndRespond(w, "Missing user name", nil)
		return
	}

	byteData, err := io.ReadAll(request.Body)
	if err != nil {
		logAndRespond(w, "Error parsing request", err)
		return
	}

	var parsedData map[string]string
	if err = json.Unmarshal(byteData, &parsedData); err != nil {
		logAndRespond(w, "Error parsing json", err)
		return
	}

	confirmationCode := parsedData["confirmationCode"]
	if confirmationCode == "" {
		logAndRespond(w, "Missing confirmation code", nil)
		return
	}

	output, err := handler.CognitoClient.ConfirmSignUp(context.TODO(), &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(config.CognitoClientID()),
		Username:         aws.String(userEmail),
		ConfirmationCode: aws.String(confirmationCode),
		SecretHash:       aws.String(config.CognitoSecretHash(userEmail)),
	})
	if err != nil {

		logAndRespond(w, "Error confirming user", err)
		return
	}

	helpers.SuccessResponse(w, &output)
}

// SignIn signs up a user with Amazon Cognito.
func (handler Handler) SignIn(w http.ResponseWriter, request *http.Request) {
	byteData, err := io.ReadAll(request.Body)
	if err != nil {
		logAndRespond(w, "Error parsing request", err)
		return
	}

	var parsedData map[string]string
	if err = json.Unmarshal(byteData, &parsedData); err != nil {
		logAndRespond(w, "Error parsing json", err)
		return
	}

	userEmail, password := parsedData["userEmail"], parsedData["password"]
	if userEmail == "" {
		logAndRespond(w, "Missing user name", nil)
		return
	}

	if password == "" {
		logAndRespond(w, "Missing password", nil)
		return
	}
	tmp := config.CognitoSecretHash(userEmail)
	output, err := handler.CognitoClient.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		ClientId: aws.String(config.CognitoClientID()),
		AuthParameters: map[string]string{
			"USERNAME":    userEmail,
			"PASSWORD":    password,
			"SECRET_HASH": tmp,
		},
	})
	if err != nil {
		var invalidPassword *types.NotAuthorizedException
		if errors.As(err, &invalidPassword) {
			logAndRespond(w, *invalidPassword.Message, err)
		} else {
			logAndRespond(w, "Couldn't sign in user", err)
		}

		return
	}

	response := map[string]string{
		"AccessToken":  *output.AuthenticationResult.AccessToken,
		"IdToken":      *output.AuthenticationResult.IdToken,
		"RefreshToken": *output.AuthenticationResult.RefreshToken,
		"ExpiresIn":    strconv.Itoa(int(output.AuthenticationResult.ExpiresIn)),
	}

	helpers.SuccessResponse(w, &response)
}

// Refresh signs up a user with Amazon Cognito.
func (handler Handler) Refresh(w http.ResponseWriter, request *http.Request) {
	// Get the JWT from the Authorization header
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		helpers.BadRequest(w, "Missing Authorization Header")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		helpers.BadRequest(w, "Invalid Authorization header format")
		return
	}

	token, err := helpers.VerifyToken(tokenString, false)
	if err != nil {
		helpers.BadRequest(w, err)
		return
	}

	byteData, err := io.ReadAll(request.Body)
	if err != nil {
		logAndRespond(w, "Error parsing request", err)
		return
	}

	var parsedData map[string]string
	if err = json.Unmarshal(byteData, &parsedData); err != nil {
		logAndRespond(w, "Error parsing json", err)
		return
	}

	refreshToken := parsedData["refreshToken"]
	if refreshToken == "" {
		logAndRespond(w, "Missing refresh token", nil)
		return
	}

	claims, err := helpers.ExtractClaims(token)

	output, err := handler.CognitoClient.InitiateAuth(context.TODO(), &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "REFRESH_TOKEN_AUTH",
		ClientId: aws.String(config.CognitoClientID()),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshToken,
			"SECRET_HASH":   config.CognitoSecretHash(claims["username"].(string)),
		},
	})
	if err != nil {

		logAndRespond(w, "Couldn't sign in user", err)

		return
	}

	response := map[string]string{
		"AccessToken": *output.AuthenticationResult.AccessToken,
		"IdToken":     *output.AuthenticationResult.IdToken,
		"ExpiresIn":   strconv.Itoa(int(output.AuthenticationResult.ExpiresIn)),
	}

	helpers.SuccessResponse(w, &response)
}
