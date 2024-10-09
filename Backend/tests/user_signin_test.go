// Package tests is where all of out unit tests are described.
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"organize-this/config"
	"organize-this/controllers"
	"organize-this/repository"
	"organize-this/tests/mocks"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
)

type userSignInResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type userSignInTestCase struct {
	testName          string
	validData         bool
	incorrectPassword bool
	expectedHTTP      int
	expectedMessage   string
	expectedBody      string
	UserEmail         string `json:"userEmail"`
	Password          string `json:"password"`
}

var userSignInEndpoint = "/v1/token"

func setupUserSignInTest(t *testing.T) (*http.Client, *httptest.Server, *mocks.MockCognitoClient) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	postgres, _ := mocks.NewMockDB()
	redis, _ := redismock.NewClientMock()
	cognito := mocks.NewMockCognitoClient(ctrl)
	handler := controllers.Handler{
		Repository:    &repository.Repository{Database: postgres, Cache: redis},
		CognitoClient: cognito,
		TokenHelper:   mocks.NewMockTokenHelper(ctrl),
	}

	r := chi.NewRouter()
	r.Post(userSignInEndpoint, handler.SignIn)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &http.Client{}, srv, cognito
}

func setupUserSignInMockExpectations(cognito *mocks.MockCognitoClient, userEmail string, password string, incorrectPassword bool) {

	if incorrectPassword {
		cognito.EXPECT().InitiateAuth(gomock.All(), &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: "USER_PASSWORD_AUTH",
			ClientId: aws.String(config.CognitoClientID()),
			AuthParameters: map[string]string{
				"USERNAME":    userEmail,
				"PASSWORD":    password,
				"SECRET_HASH": config.CognitoSecretHash(userEmail),
			},
		}).Return(nil,
			&types.NotAuthorizedException{
				Message: aws.String("Password or username incorrect, please try again."),
			},
		)
	} else {
		cognito.EXPECT().InitiateAuth(gomock.All(), &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: "USER_PASSWORD_AUTH",
			ClientId: aws.String(config.CognitoClientID()),
			AuthParameters: map[string]string{
				"USERNAME":    userEmail,
				"PASSWORD":    password,
				"SECRET_HASH": config.CognitoSecretHash(userEmail),
			},
		}).Return(&cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &types.AuthenticationResultType{
				AccessToken:  aws.String("123"),
				IdToken:      aws.String("123"),
				RefreshToken: aws.String("123"),
				ExpiresIn:    3600,
			},
		}, nil)
	}
}

func validateUserSignInResponse(t *testing.T, res *http.Response, expectedHTTP int, expectedMessage string, expectedData string) {
	if res.StatusCode != expectedHTTP {
		t.Errorf("Expected status code to be: %d. Got: %d.", expectedHTTP, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	contents := userSignInResponse{}
	err = json.Unmarshal(data, &contents)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	if contents.Message != expectedMessage {
		t.Errorf("Expected message to be %s. Got: %s", expectedMessage, contents.Message)
	}

	if expectedData != "" && contents.Data != expectedData {
		t.Errorf("Expected body to be %s. Got: %s", expectedData, contents.Data)
	}
}

// TestUserSignIn runs the unit tests for the UserSignIn function.
func TestUserSignIn(t *testing.T) {
	cases := []userSignInTestCase{
		{
			testName:          "BEUT-43: User Sign In Valid Data",
			validData:         true,
			incorrectPassword: false,
			expectedHTTP:      http.StatusOK,
			expectedMessage:   "success",
			expectedBody:      "",
			UserEmail:         "test@test.gmail.com",
			Password:          "password",
		},
		{
			testName:          "BEUT-44: User Sign In Incorrect Password",
			validData:         true,
			incorrectPassword: true,
			expectedHTTP:      http.StatusBadRequest,
			expectedMessage:   "data validation failed",
			expectedBody:      "Password or username incorrect, please try again.",
			UserEmail:         "test@test.gmail.com",
			Password:          "password",
		},
		{
			testName:        "BEUT-45: User Sign In Missing Email",
			validData:       false,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Missing user email",
			UserEmail:       "",
			Password:        "152634",
		},
		{
			testName:        "BEUT-46: User Sign In Missing Email",
			validData:       false,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Missing password",
			UserEmail:       "test@test.test",
			Password:        "",
		},
	}

	for _, tc := range cases {
		client, srv, cognito := setupUserSignInTest(t)
		t.Run(tc.testName, func(t *testing.T) {

			if tc.validData {
				setupUserSignInMockExpectations(cognito, tc.UserEmail, tc.Password, tc.incorrectPassword)
			}

			url := fmt.Sprintf("%s%s", srv.URL, userSignInEndpoint)

			jsonBody, err := json.Marshal(tc)
			if err != nil {
				t.Fatalf("Error marshalling json: %v", err)
			}

			res, err := client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer res.Body.Close()

			validateUserSignInResponse(t, res, tc.expectedHTTP, tc.expectedMessage, tc.expectedBody)

		})
	}
}
