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
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
)

type userRefreshResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type userRefreshTestCase struct {
	testName        string
	validData       bool
	badRequest      bool
	expectedHTTP    int
	expectedMessage string
	expectedBody    string
	userName        string
	RefreshToken    string `json:"refreshToken"`
	IDToken         string `json:"idToken"`
}

var userRefreshEndpoint = "/v1/token"

func setupUserRefreshTest(t *testing.T) (*http.Client, *httptest.Server, *mocks.MockCognitoClient, *mocks.MockTokenHelper) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	postgres, _ := mocks.NewMockDB()
	redis, _ := redismock.NewClientMock()
	cognito := mocks.NewMockCognitoClient(ctrl)
	tokenHelper := mocks.NewMockTokenHelper(ctrl)
	handler := controllers.Handler{
		Repository:    &repository.Repository{Database: postgres, Cache: redis},
		CognitoClient: cognito,
		TokenHelper:   tokenHelper,
	}

	r := chi.NewRouter()
	r.Put(userRefreshEndpoint, handler.Refresh)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &http.Client{}, srv, cognito, tokenHelper
}

func setupUserRefreshMockExpectations(cognito *mocks.MockCognitoClient, tokenHelper *mocks.MockTokenHelper, refreshToken string, userName string, badRequest bool) {

	tokenHelper.EXPECT().VerifyToken(gomock.Any(), false).Return(&jwt.Token{}, nil)
	tokenHelper.EXPECT().ExtractClaims(gomock.Any()).Return(jwt.MapClaims{"cognito:username": userName}, nil)

	if badRequest {
		cognito.EXPECT().InitiateAuth(gomock.All(), &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: "REFRESH_TOKEN_AUTH",
			ClientId: aws.String(config.CognitoClientID()),
			AuthParameters: map[string]string{
				"REFRESH_TOKEN": refreshToken,
				"SECRET_HASH":   config.CognitoSecretHash(userName),
			},
		}).Return(nil, &types.NotAuthorizedException{
			Message: aws.String("Incorrect token or userName."),
		})
	} else {
		cognito.EXPECT().InitiateAuth(gomock.All(), &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: "REFRESH_TOKEN_AUTH",
			ClientId: aws.String(config.CognitoClientID()),
			AuthParameters: map[string]string{
				"REFRESH_TOKEN": refreshToken,
				"SECRET_HASH":   config.CognitoSecretHash(userName),
			},
		}).Return(&cognitoidentityprovider.InitiateAuthOutput{
			AuthenticationResult: &types.AuthenticationResultType{
				AccessToken: aws.String("123"),
				IdToken:     aws.String("123"),
				ExpiresIn:   3600,
			},
		}, nil)
	}

}

func validateUserRefreshResponse(t *testing.T, res *http.Response, expectedHTTP int, expectedMessage string, expectedData string) {
	if res.StatusCode != expectedHTTP {
		t.Errorf("Expected status code to be: %d. Got: %d.", expectedHTTP, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	contents := userRefreshResponse{}
	err = json.Unmarshal(data, &contents)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	if contents.Message != expectedMessage {
		t.Errorf("Expected message to be %s. Got: %s : %s", expectedMessage, contents.Message, contents.Data)
	}

	if expectedData != "" && contents.Data != expectedData {
		t.Errorf("Expected body to be %s. Got: %s", expectedData, contents.Data)
	}
}

// TestUserRefresh runs the unit tests for the Refresh function.
func TestUserRefresh(t *testing.T) {
	cases := []userRefreshTestCase{
		{
			testName:        "BEUT-50: User Refresh Valid Data",
			validData:       true,
			badRequest:      false,
			expectedHTTP:    http.StatusOK,
			expectedMessage: "success",
			expectedBody:    "",
			userName:        "testuser",
			RefreshToken:    "lkjanbdljkahvf",
			IDToken:         "eyJraWQiOiJQNW",
		},
		{
			testName:        "BEUT-51: User Refresh Bad Data",
			validData:       true,
			badRequest:      true,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Couldn't refresh user",
			userName:        "testuser1",
			RefreshToken:    "lkjanbdljkahvf",
			IDToken:         "eyJraWQiOiJQNW",
		},
		{
			testName:        "BEUT-52: User Refresh Missing Refresh Token",
			validData:       false,
			badRequest:      false,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Missing refresh token",
			userName:        "testuser2",
			RefreshToken:    "",
			IDToken:         "eyJraWQiOiJQNW",
		},
		{
			testName:        "BEUT-53: User Refresh Missing ID Token",
			validData:       false,
			badRequest:      false,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Missing id token",
			userName:        "testuser3",
			RefreshToken:    "lkjanbdljkahvf",
			IDToken:         "",
		},
	}

	for _, tc := range cases {
		client, srv, cognito, tokenHelper := setupUserRefreshTest(t)
		t.Run(tc.testName, func(t *testing.T) {

			if tc.validData {
				setupUserRefreshMockExpectations(cognito, tokenHelper, tc.RefreshToken, tc.userName, tc.badRequest)
			}

			url := fmt.Sprintf("%s%s", srv.URL, userRefreshEndpoint)

			jsonBody, err := json.Marshal(tc)
			if err != nil {
				t.Fatalf("Error marshalling json: %v", err)
			}

			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to build request: %v", err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer res.Body.Close()

			validateUserRefreshResponse(t, res, tc.expectedHTTP, tc.expectedMessage, tc.expectedBody)

		})
	}
}
