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

type userSignUpResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type userSignUpTestCase struct {
	testName        string
	validData       bool
	badPassword     bool
	expectedHTTP    int
	expectedMessage string
	expectedBody    string
	UserEmail       string `json:"userEmail"`
	Password        string `json:"password"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Birthday        string `json:"birthday"`
}

var userSignUpEndpoint = "/v1/user"

func setupUserSignUpTest(t *testing.T) (*http.Client, *httptest.Server, *mocks.MockCognitoClient) {
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
	r.Post(userSignUpEndpoint, handler.SignUp)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &http.Client{}, srv, cognito
}

func setupUserSignUpMockExpectations(cognito *mocks.MockCognitoClient, userEmail string, password string, firstName string, lastName string, birthday string, baddpassword bool) {
	if baddpassword {
		cognito.EXPECT().SignUp(gomock.All(), &cognitoidentityprovider.SignUpInput{
			ClientId: aws.String(config.CognitoClientID()),
			Password: aws.String(password),
			Username: aws.String(userEmail),
			UserAttributes: []types.AttributeType{
				{Name: aws.String("given_name"), Value: aws.String(firstName)},
				{Name: aws.String("family_name"), Value: aws.String(lastName)},
				{Name: aws.String("birthdate"), Value: aws.String(birthday)},
			},
			SecretHash: aws.String(config.CognitoSecretHash(userEmail)),
		}).Return(nil,
			&types.InvalidPasswordException{
				Message: aws.String("Invalid password configuration, please try again."),
			},
		)
	} else {
		cognito.EXPECT().SignUp(gomock.All(), &cognitoidentityprovider.SignUpInput{
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
	}
}

func validateUserSignUpResponse(t *testing.T, res *http.Response, expectedHTTP int, expectedMessage string, expectedData string) {
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

// TestCreateEntity runs the unit tests for invalid cases.
func TestUserSignUp(t *testing.T) {
	cases := []userSignUpTestCase{
		{
			testName:        "BEUT-32: User Sign Up Valid Data",
			validData:       true,
			expectedHTTP:    http.StatusOK,
			expectedMessage: "success",
			expectedBody:    "",
			UserEmail:       "test@test.gmail.com",
			Password:        "password",
			FirstName:       "test",
			LastName:        "test",
			Birthday:        "06/06/06",
		},
		{
			testName:        "BEUT-33: User Sign Up Invalid Password",
			validData:       true,
			badPassword:     true,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Invalid password configuration, please try again.",
			UserEmail:       "test@test.gmail.com",
			Password:        "password",
			FirstName:       "test",
			LastName:        "test",
			Birthday:        "06/06/06",
		},
		{
			testName:        "BEUT-34: User Sign Up Missing Email",
			validData:       false,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Missing user name",
			UserEmail:       "",
			Password:        "password",
			FirstName:       "test",
			LastName:        "test",
			Birthday:        "06/06/06",
		},
		{
			testName:        "BEUT-35: User Sign Up Missing Password",
			validData:       false,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Missing password",
			UserEmail:       "test@test.gmail.com",
			Password:        "",
			FirstName:       "test",
			LastName:        "test",
			Birthday:        "06/06/06",
		},
		{
			testName:        "BEUT-36: User Sign Up Missing First Name",
			validData:       false,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Missing first name",
			UserEmail:       "test@test.test",
			Password:        "password",
			FirstName:       "",
			LastName:        "test",
			Birthday:        "06/06/06",
		},
		{
			testName:        "BEUT-37: User Sign Up Missing Last Name",
			validData:       false,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Missing last name",
			UserEmail:       "test@test.test",
			Password:        "password",
			FirstName:       "test",
			LastName:        "",
			Birthday:        "06/06/06",
		},
		{
			testName:        "BEUT-38: User Sign Up Missing Birthday",
			validData:       false,
			expectedHTTP:    http.StatusBadRequest,
			expectedMessage: "data validation failed",
			expectedBody:    "Missing birthday",
			UserEmail:       "test@test.test",
			Password:        "password",
			FirstName:       "test",
			LastName:        "test",
			Birthday:        "",
		},
	}

	for _, tc := range cases {
		client, srv, cognito := setupUserSignUpTest(t)
		t.Run(tc.testName, func(t *testing.T) {

			if tc.validData {
				setupUserSignUpMockExpectations(cognito, tc.UserEmail, tc.Password, tc.FirstName, tc.LastName, tc.Birthday, tc.badPassword)
			}

			url := fmt.Sprintf("%s%s", srv.URL, userSignUpEndpoint)

			jsonBody, err := json.Marshal(tc)
			if err != nil {
				t.Fatalf("Error marshalling json: %v", err)
			}

			res, err := client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer res.Body.Close()

			validateUserSignUpResponse(t, res, tc.expectedHTTP, tc.expectedMessage, tc.expectedBody)

		})
	}
}
