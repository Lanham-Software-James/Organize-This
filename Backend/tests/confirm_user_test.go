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

type confirmUserResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type confirmUserTestCase struct {
	testName         string
	validData        bool
	incorrectCode    bool
	expectedHTTP     int
	expectedMessage  string
	expectedBody     string
	UserEmail        string `json:"userEmail"`
	ConfirmationCode string `json:"confirmationCode"`
}

var confirmUserEndpoint = "/v1/user"

func setupConfirmUserTest(t *testing.T) (*http.Client, *httptest.Server, *mocks.MockCognitoClient) {
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
	r.Put(confirmUserEndpoint, handler.ConfirmSignUp)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &http.Client{}, srv, cognito
}

func setupConfirmUserMockExpectations(cognito *mocks.MockCognitoClient, userEmail string, confirmationCode string, incorrectCode bool) {

	if incorrectCode {
		cognito.EXPECT().ConfirmSignUp(gomock.All(), &cognitoidentityprovider.ConfirmSignUpInput{
			ClientId:         aws.String(config.CognitoClientID()),
			Username:         aws.String(userEmail),
			ConfirmationCode: aws.String(confirmationCode),
			SecretHash:       aws.String(config.CognitoSecretHash(userEmail)),
		}).Return(nil,
			&types.CodeMismatchException{
				Message: aws.String("Invalid verification code provided, please try again."),
			},
		)
	} else {
		cognito.EXPECT().ConfirmSignUp(gomock.All(), &cognitoidentityprovider.ConfirmSignUpInput{
			ClientId:         aws.String(config.CognitoClientID()),
			Username:         aws.String(userEmail),
			ConfirmationCode: aws.String(confirmationCode),
			SecretHash:       aws.String(config.CognitoSecretHash(userEmail)),
		})
	}

}

func validateConfirmUserResponse(t *testing.T, res *http.Response, expectedHTTP int, expectedMessage string, expectedData string) {
	if res.StatusCode != expectedHTTP {
		t.Errorf("Expected status code to be: %d. Got: %d.", expectedHTTP, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	contents := userSignUpResponse{}
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

// TestConfirmUser runs the unit tests for the ConfirmUser function.
func TestConfirmUser(t *testing.T) {
	cases := []confirmUserTestCase{
		{
			testName:         "BEUT-40: Confirm User Valid Data",
			validData:        true,
			incorrectCode:    false,
			expectedHTTP:     http.StatusOK,
			expectedMessage:  "success",
			expectedBody:     "",
			UserEmail:        "test@test.gmail.com",
			ConfirmationCode: "XXXXXXX",
		},
		{
			testName:         "BEUT-41: Confirm User Incorrect Code",
			validData:        true,
			incorrectCode:    true,
			expectedHTTP:     http.StatusBadRequest,
			expectedMessage:  "data validation failed",
			expectedBody:     "Incorrect confirmation code",
			UserEmail:        "test@test.gmail.com",
			ConfirmationCode: "XXXXXXX",
		},
		{
			testName:         "BEUT-42: Confirm User Missing Email",
			validData:        false,
			expectedHTTP:     http.StatusBadRequest,
			expectedMessage:  "data validation failed",
			expectedBody:     "Missing user email",
			UserEmail:        "",
			ConfirmationCode: "152634",
		},
		{
			testName:         "BEUT-43: Confirm User Missing Email",
			validData:        false,
			expectedHTTP:     http.StatusBadRequest,
			expectedMessage:  "data validation failed",
			expectedBody:     "Missing confirmation code",
			UserEmail:        "test@test.test",
			ConfirmationCode: "",
		},
	}

	for _, tc := range cases {
		client, srv, cognito := setupConfirmUserTest(t)
		t.Run(tc.testName, func(t *testing.T) {

			if tc.validData {
				setupConfirmUserMockExpectations(cognito, tc.UserEmail, tc.ConfirmationCode, tc.incorrectCode)
			}

			url := fmt.Sprintf("%s%s", srv.URL, userSignUpEndpoint)

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

			validateConfirmUserResponse(t, res, tc.expectedHTTP, tc.expectedMessage, tc.expectedBody)

		})
	}
}
