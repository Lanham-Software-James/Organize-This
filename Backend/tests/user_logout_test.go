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
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
)

type userLogOutResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type userLogOutTestCase struct {
	testName        string
	validData       bool
	expectedHTTP    int
	expectedMessage string
	expectedBody    string
	RefreshToken    string `json:"refreshToken"`
}

var userLogOutEndpoint = "/v1/token"

func setupUserLogOutTest(t *testing.T) (*http.Client, *httptest.Server, *mocks.MockCognitoClient) {
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
	r.Delete(userLogOutEndpoint, handler.LogOut)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &http.Client{}, srv, cognito
}

func setupUserLogOutMockExpectations(cognito *mocks.MockCognitoClient, refreshToken string) {

	cognito.EXPECT().RevokeToken(gomock.All(), &cognitoidentityprovider.RevokeTokenInput{
		ClientId:     aws.String(config.CognitoClientID()),
		ClientSecret: aws.String(config.CognitoClientSecret()),
		Token:        aws.String(refreshToken),
	})
}

func validateUserLogOutResponse(t *testing.T, res *http.Response, expectedHTTP int, expectedMessage string, expectedData string) {
	if res.StatusCode != expectedHTTP {
		t.Errorf("Expected status code to be: %d. Got: %d.", expectedHTTP, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	contents := userLogOutResponse{}
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

// TestUserLogOut runs the unit tests for the UserLogOut function.
func TestUserLogOut(t *testing.T) {
	cases := []userLogOutTestCase{
		{
			testName:        "BEUT-47: User Log Out Valid Data",
			validData:       true,
			expectedHTTP:    http.StatusOK,
			expectedMessage: "success",
			expectedBody:    "",
			RefreshToken:    "jkadlkfjhadsfakbcvadsuiofha",
		},
		{
			testName:        "BEUT-48: User Log Out Invalid Data",
			validData:       false,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Missing refresh token",
			RefreshToken:    "",
		},
	}

	for _, tc := range cases {
		client, srv, cognito := setupUserLogOutTest(t)
		t.Run(tc.testName, func(t *testing.T) {

			if tc.validData {
				setupUserLogOutMockExpectations(cognito, tc.RefreshToken)
			}

			url := fmt.Sprintf("%s%s", srv.URL, userLogOutEndpoint)

			jsonBody, err := json.Marshal(tc)
			if err != nil {
				t.Fatalf("Error marshalling json: %v", err)
			}

			req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to build request: %v", err)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer res.Body.Close()

			validateUserLogOutResponse(t, res, tc.expectedHTTP, tc.expectedMessage, tc.expectedBody)

		})
	}
}
