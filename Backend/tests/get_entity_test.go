// Package tests is where all of out unit tests are described.
package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"organize-this/controllers"
	"organize-this/repository"
	"organize-this/tests/mocks"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
)

type getEntitySingleResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type getEntityTestCase struct {
	testName       string
	testUser       string
	testValidInput bool
	category       string
	id             string
}

var getEntityEndpoint = "/v1/entity"
var getEntityParameters = "/{category}/{id}"

func setupGetEntityTest(t *testing.T, userName string) (*http.Client, *httptest.Server, sqlmock.Sqlmock) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	postgres, mockDB := mocks.NewMockDB()
	redis, _ := redismock.NewClientMock()
	handler := controllers.Handler{
		Repository:    &repository.Repository{Database: postgres, Cache: redis},
		CognitoClient: mocks.NewMockCognitoClient(ctrl),
		TokenHelper:   mocks.NewMockTokenHelper(ctrl),
	}

	r := chi.NewRouter()
	r.Use(mocks.MockJWTMiddleware(userName))
	r.Get(getEntityEndpoint+getEntityParameters, handler.GetEntity)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &http.Client{}, srv, mockDB
}

func setupGetEntityMockExpectations(mockDB *sqlmock.Sqlmock, category string, userName string, entityID string) {
	tableName := category + "s"
	if category == "shelf" {
		tableName = "shelves"
	}

	entityIDInt, _ := strconv.ParseInt(entityID, 10, 64)

	expectedMainSQL := fmt.Sprintf(`SELECT \* FROM "%s" WHERE user_id = \$1 AND "%s"."deleted_at" IS NULL AND "%s"."id" = \$2 ORDER BY "%s"."id" LIMIT 1`, tableName, tableName, tableName, tableName)

	rows := sqlmock.NewRows([]string{"id", "name", "notes", "user_id", "created_at", "updated_at", "deleted_at", "parent_id", "parent_category"}).
		AddRow(entityID, "Entity 1", "Notes", userName, time.Now(), time.Now(), nil, 0, "")

	(*mockDB).ExpectQuery(expectedMainSQL).
		WithArgs(userName, entityIDInt).
		WillReturnRows(rows)
}

func validateGetEntitySuccessResponse(t *testing.T, res *http.Response, mockDB sqlmock.Sqlmock) {
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	contents := getEntitySingleResponse{}
	err = json.Unmarshal(data, &contents)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	if contents.Message != "success" {
		t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
	}

	dataType := reflect.TypeOf(contents.Data).String()
	if dataType != "map[string]interface {}" {
		t.Errorf("Expected data to be type models.GetEntityResponse. Got: %v", dataType)
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("PostGres expectations were not met: %v", err)
	}
}

// TestCreateEntity runs the unit tests for invalid cases.
func TestGetEntity(t *testing.T) {
	cases := []getEntityTestCase{
		{
			testName:       "BEUT-90: Get Entity Invalid Category",
			testUser:       "testuser0",
			testValidInput: false,
			category:       "test",
			id:             "1",
		},
		{
			testName:       "BEUT-91: Get Entity Invalid ID",
			testUser:       "testuser0",
			testValidInput: false,
			category:       "item",
			id:             "String",
		},
		{
			testName:       "BEUT-92: Get Entity Item",
			testUser:       "testuser0",
			testValidInput: true,
			category:       "item",
			id:             "1",
		},
		{
			testName:       "BEUT-93: Get Entity Container",
			testUser:       "testuser0",
			testValidInput: true,
			category:       "container",
			id:             "1",
		},
		{
			testName:       "BEUT-94: Get Entity Shelf",
			testUser:       "testuser0",
			testValidInput: true,
			category:       "shelf",
			id:             "1",
		},
		{
			testName:       "BEUT-95: Get Entity Shelving Unit",
			testUser:       "testuser0",
			testValidInput: true,
			category:       "shelving_unit",
			id:             "1",
		},
		{
			testName:       "BEUT-96: Get Entity Room",
			testUser:       "testuser0",
			testValidInput: true,
			category:       "room",
			id:             "1",
		},
		{
			testName:       "BEUT-97: Get Entity Building",
			testUser:       "testuser0",
			testValidInput: true,
			category:       "building",
			id:             "1",
		},
	}

	for _, tc := range cases {
		client, srv, mockDB := setupGetEntityTest(t, tc.testUser)
		t.Run(tc.testName, func(t *testing.T) {

			if tc.testValidInput {
				setupGetEntityMockExpectations(&mockDB, tc.category, tc.testUser, tc.id)
			}

			res, err := client.Get(fmt.Sprintf("%s%s/%s/%s", srv.URL, getEntityEndpoint, tc.category, tc.id))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer res.Body.Close()

			if tc.testValidInput {
				validateGetEntitySuccessResponse(t, res, mockDB)
			} else if (!tc.testValidInput) && (res.StatusCode != http.StatusBadRequest) {
				t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
			}

		})
	}
}
