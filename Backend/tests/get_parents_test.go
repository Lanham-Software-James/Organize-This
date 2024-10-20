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
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
)

type getParentsSingleResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type getParentsTestCase struct {
	testName       string
	testUser       string
	testValidInput bool
	testCacheHit   bool
	category       string
}

var getParentsEndpoint = "/v1/parents"
var getParentsParameters = "/{category}"

func setupGetParentsTest(t *testing.T, userName string) (*http.Client, *httptest.Server, sqlmock.Sqlmock, redismock.ClientMock) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	postgres, mockDB := mocks.NewMockDB()
	redis, mockCache := redismock.NewClientMock()
	handler := controllers.Handler{
		Repository:    &repository.Repository{Database: postgres, Cache: redis},
		CognitoClient: mocks.NewMockCognitoClient(ctrl),
		TokenHelper:   mocks.NewMockTokenHelper(ctrl),
	}

	r := chi.NewRouter()
	r.Use(mocks.MockJWTMiddleware(userName))
	r.Get(getParentsEndpoint+getParentsParameters, handler.GetParents)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &http.Client{}, srv, mockDB, mockCache
}

func setupGetParentsCacheMissMockExpectations(mockDB *sqlmock.Sqlmock, mockCache redismock.ClientMock, category string, userName string) {
	switch category {
	case "item":
		cacheKey := fmt.Sprintf(`{"User":"%s","Function":"GetItemParents"}`, userName)
		expectedSQL := `
        (SELECT 'room' AS category, id, name FROM rooms WHERE user_id = $1)
        UNION ALL
        (SELECT 'shelf' AS category, id, name FROM shelves WHERE user_id = $2)
        UNION ALL
        (SELECT 'container' AS category, id, name FROM containers WHERE user_id = $3)`

		(*mockDB).ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(userName, userName, userName).
			WillReturnRows(sqlmock.NewRows([]string{"category", "id", "name"}).
				AddRow("room", 1, "Room 1").
				AddRow("shelf", 2, "Shelf 1").
				AddRow("container", 3, "Container 1"))

		mockCache.ExpectGet(cacheKey).RedisNil()
		mockCache.Regexp().ExpectSet(cacheKey, ".*", 5*time.Minute).SetVal("OK")
		break
	case "container":
		cacheKey := fmt.Sprintf(`{"User":"%s","Function":"GetContainerParents"}`, userName)
		expectedSQL := `
        (SELECT 'room' AS category, id, name FROM rooms WHERE user_id = $1)
        UNION ALL
        (SELECT 'shelf' AS category, id, name FROM shelves WHERE user_id = $2)`

		(*mockDB).ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(userName, userName).
			WillReturnRows(sqlmock.NewRows([]string{"category", "id", "name"}).
				AddRow("room", 1, "Room 1").
				AddRow("shelf", 2, "Shelf 1"))

		mockCache.ExpectGet(cacheKey).RedisNil()
		mockCache.Regexp().ExpectSet(cacheKey, ".*", 5*time.Minute).SetVal("OK")
		break
	case "shelf":
		cacheKey := fmt.Sprintf(`{"User":"%s","Function":"GetShelfParents"}`, userName)
		expectedSQL := `
        (SELECT 'shelving_unit' AS category, id, name FROM shelving_units WHERE user_id = $1)`

		(*mockDB).ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(userName).
			WillReturnRows(sqlmock.NewRows([]string{"category", "id", "name"}).
				AddRow("shelving_unit", 1, "Shelving Unit 1"))

		mockCache.ExpectGet(cacheKey).RedisNil()
		mockCache.Regexp().ExpectSet(cacheKey, ".*", 5*time.Minute).SetVal("OK")
		break
	case "shelving_unit":
		cacheKey := fmt.Sprintf(`{"User":"%s","Function":"GetShelving_unitParents"}`, userName)
		expectedSQL := `
        (SELECT 'room' AS category, id, name FROM rooms WHERE user_id = $1)`

		(*mockDB).ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(userName).
			WillReturnRows(sqlmock.NewRows([]string{"category", "id", "name"}).
				AddRow("room", 1, "Room 1"))

		mockCache.ExpectGet(cacheKey).RedisNil()
		mockCache.Regexp().ExpectSet(cacheKey, ".*", 5*time.Minute).SetVal("OK")
		break
	case "room":
		cacheKey := fmt.Sprintf(`{"User":"%s","Function":"GetRoomParents"}`, userName)
		expectedSQL := `
        (SELECT 'building' AS category, id, name FROM buildings WHERE user_id = $1)`

		(*mockDB).ExpectQuery(regexp.QuoteMeta(expectedSQL)).
			WithArgs(userName).
			WillReturnRows(sqlmock.NewRows([]string{"category", "id", "name"}).
				AddRow("building", 1, "Building 1"))

		mockCache.ExpectGet(cacheKey).RedisNil()
		mockCache.Regexp().ExpectSet(cacheKey, ".*", 5*time.Minute).SetVal("OK")
		break
	}

}

func setupGetParentsCacheHitMockExpectations(mockCache redismock.ClientMock, category string, userName string) {
	switch category {
	case "item":
		cacheKey := fmt.Sprintf(`{"User":"%s","Function":"GetItemParents"}`, userName)
		mockCache.ExpectGet(cacheKey).SetVal(`[
			{"ID":13,"Name":"Cat Room","Category":"room"},
			{"ID":11,"Name":"Shelf 1","Category":"shelf"},
			{"ID":12,"Name":"Basket","Category":"container"}
		]`)
		break
	case "container":
		cacheKey := fmt.Sprintf(`{"User":"%s","Function":"GetContainerParents"}`, userName)
		mockCache.ExpectGet(cacheKey).SetVal(`[
			{"ID":13,"Name":"Cat Room","Category":"room"},
			{"ID":11,"Name":"Shelf 1","Category":"shelf"}
		]`)
		break
	case "shelf":
		cacheKey := fmt.Sprintf(`{"User":"%s","Function":"GetShelfParents"}`, userName)
		mockCache.ExpectGet(cacheKey).SetVal(`[
			{"ID":11,"Name":"Shelving Unit 1","Category":"shelving_unit"}
		]`)
		break
	case "shelving_unit":
		cacheKey := fmt.Sprintf(`{"User":"%s","Function":"GetShelving_unitParents"}`, userName)
		mockCache.ExpectGet(cacheKey).SetVal(`[
			{"ID":11,"Name":"Room 1","Category":"room"}
		]`)
		break
	case "room":
		cacheKey := fmt.Sprintf(`{"User":"%s","Function":"GetRoomParents"}`, userName)
		mockCache.ExpectGet(cacheKey).SetVal(`[
			{"ID":11,"Name":"Building 1","Category":"building"}
		]`)
		break
	}

}

func validateGetParentsSuccessResponse(t *testing.T, res *http.Response, mockDB sqlmock.Sqlmock, mockCache redismock.ClientMock) {
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	contents := getParentsSingleResponse{}
	err = json.Unmarshal(data, &contents)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	if contents.Message != "success" {
		t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
	}

	dataType := reflect.TypeOf(contents.Data).String()
	if dataType != "[]interface {}" {
		t.Errorf("Expected data to be type []interface{}. Got: %v", dataType)
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("PostGres expectations were not met: %v", err)
	}

	if err := mockCache.ExpectationsWereMet(); err != nil {
		t.Errorf("Redis expectations were not met: %v", err)
	}
}

// TestCreateParents runs the unit tests for invalid cases.
func TestGetParents(t *testing.T) {
	cases := []getParentsTestCase{
		{
			testName:       "BEUT-98: Get Parents Invalid Category",
			testUser:       "testuser0",
			testValidInput: false,
			category:       "test",
		},
		{
			testName:       "BEUT-99: Get Parents Item Cache Miss",
			testUser:       "testuser0",
			testValidInput: true,
			testCacheHit:   false,
			category:       "item",
		},
		{
			testName:       "BEUT-100: Get Parents Item Cache Hit",
			testUser:       "testuser0",
			testValidInput: true,
			testCacheHit:   true,
			category:       "item",
		},
		{
			testName:       "BEUT-101: Get Parents Container Cache Miss",
			testUser:       "testuser0",
			testValidInput: true,
			testCacheHit:   false,
			category:       "container",
		},
		{
			testName:       "BEUT-102: Get Parents Container Cache Hit",
			testUser:       "testuser0",
			testValidInput: true,
			testCacheHit:   true,
			category:       "container",
		},
		{
			testName:       "BEUT-103: Get Parents Shelf Cache Miss",
			testUser:       "testuser0",
			testValidInput: true,
			testCacheHit:   false,
			category:       "shelf",
		},
		{
			testName:       "BEUT-104: Get Parents Shelf Cache Hit",
			testUser:       "testuser0",
			testValidInput: true,
			testCacheHit:   true,
			category:       "shelf",
		},
		{
			testName:       "BEUT-105: Get Parents Shelving Unit Cache Miss",
			testUser:       "testuser0",
			testValidInput: true,
			testCacheHit:   false,
			category:       "shelving_unit",
		},
		{
			testName:       "BEUT-106: Get Parents Shelving Unit Cache Hit",
			testUser:       "testuser0",
			testValidInput: true,
			testCacheHit:   true,
			category:       "shelving_unit",
		},
		{
			testName:       "BEUT-107: Get Parents Room Cache Miss",
			testUser:       "testuser0",
			testValidInput: true,
			testCacheHit:   false,
			category:       "room",
		},
		{
			testName:       "BEUT-108: Get Parents Room Cache Hit",
			testUser:       "testuser0",
			testValidInput: true,
			testCacheHit:   true,
			category:       "room",
		},
	}

	for _, tc := range cases {
		client, srv, mockDB, mockCache := setupGetParentsTest(t, tc.testUser)
		t.Run(tc.testName, func(t *testing.T) {

			if tc.testValidInput && tc.testCacheHit {
				setupGetParentsCacheHitMockExpectations(mockCache, tc.category, tc.testUser)
			} else if tc.testValidInput {
				setupGetParentsCacheMissMockExpectations(&mockDB, mockCache, tc.category, tc.testUser)
			}

			res, err := client.Get(fmt.Sprintf("%s%s/%s", srv.URL, getParentsEndpoint, tc.category))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer res.Body.Close()

			if tc.testValidInput {
				validateGetParentsSuccessResponse(t, res, mockDB, mockCache)
			} else if (!tc.testValidInput) && (res.StatusCode != http.StatusBadRequest) {
				t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
			}

		})
	}
}
