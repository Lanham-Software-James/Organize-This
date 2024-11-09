// Package tests is where all of out unit tests are described.
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"organize-this/controllers"
	"organize-this/repository"
	"organize-this/tests/mocks"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
)

type deleteEntitySingleResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type deleteEntityTestCase struct {
	testName         string
	testUser         string
	validData        bool
	numberOfChildren int
	category         string
	id               int64
}

var deleteEntityEndpoint = "/v1/entity"
var deleteEntityParameters = "/{category}/{id}"

func setupDeleteEntityTest(t *testing.T, userName string) (*http.Client, *httptest.Server, sqlmock.Sqlmock, redismock.ClientMock) {
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
	r.Delete(deleteEntityEndpoint+deleteEntityParameters, handler.DeleteEntity)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &http.Client{}, srv, mockDB, mockCache
}

func setupDeleteEntityMockExpectations(mockDB *sqlmock.Sqlmock, mockCache redismock.ClientMock, testID int64, category string, testUser string, numberOfChildren int) {
	tableName := category + "s"
	if category == "shelf" {
		tableName = "shelves"
	}

	// Expect the SELECT operation to check if the entity exists
	selectQuery := fmt.Sprintf(`SELECT * FROM "%s" WHERE user_id = $1 AND "%s"."deleted_at" IS NULL AND "%s"."id" = $2 ORDER BY "%s"."id" LIMIT 1`, tableName, tableName, tableName, tableName)
	(*mockDB).ExpectQuery(regexp.QuoteMeta(selectQuery)).
		WithArgs(testUser, testID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "notes", "user_id", "created_at", "updated_at", "deleted_at", "parent_id", "parent_category"}).
			AddRow(testID, "Test Item", "Test Notes", testUser, time.Now(), time.Now(), nil, 0, ""))

	if category == "container" {
		// Expect the query to check for container children
		childrenQuery := `SELECT count(id) AS childrenCount FROM items WHERE user_id = $1 AND parent_id = $2 AND parent_category = $3 AND deleted_at IS NULL`
		(*mockDB).ExpectQuery(regexp.QuoteMeta(childrenQuery)).
			WithArgs(testUser, testID, category).
			WillReturnRows(sqlmock.NewRows([]string{"childrenCount"}).AddRow(numberOfChildren))
	} else if category == "shelf" {
		// Expect the query to check for shelf children
		childrenQuery := `
			SELECT
				(SELECT count(id) FROM containers WHERE user_id = $1 AND parent_id = $2 AND parent_category = $3 AND deleted_at IS NULL) +
				(SELECT count(id) FROM items WHERE user_id = $4 AND parent_id = $5 AND parent_category = $6 AND deleted_at IS NULL)
			AS childrenCount`
		(*mockDB).ExpectQuery(regexp.QuoteMeta(childrenQuery)).
			WithArgs(testUser, testID, category, testUser, testID, category).
			WillReturnRows(sqlmock.NewRows([]string{"childrenCount"}).AddRow(numberOfChildren))
	} else if category == "shelving_unit" {
		childrenQuery := `SELECT count(id) AS childrenCount FROM shelves WHERE user_id = $1 AND parent_id = $2 AND parent_category = $3 AND deleted_at IS NULL`
		(*mockDB).ExpectQuery(regexp.QuoteMeta(childrenQuery)).
			WithArgs(testUser, testID, category).
			WillReturnRows(sqlmock.NewRows([]string{"childrenCount"}).AddRow(numberOfChildren))
	} else if category == "room" {
		childrenQuery := `
		SELECT
			(SELECT count(id) FROM shelving_units WHERE user_id = $1 AND parent_id = $2 AND parent_category = $3 AND deleted_at IS NULL) +
			(SELECT count(id) FROM containers WHERE user_id = $4 AND parent_id = $5 AND parent_category = $6 AND deleted_at IS NULL) +
			(SELECT count(id) FROM items WHERE user_id = $7 AND parent_id = $8 AND parent_category = $9 AND deleted_at IS NULL)
		AS childrenCount`
		(*mockDB).ExpectQuery(regexp.QuoteMeta(childrenQuery)).
			WithArgs(testUser, testID, category, testUser, testID, category, testUser, testID, category).
			WillReturnRows(sqlmock.NewRows([]string{"childrenCount"}).AddRow(numberOfChildren))
	} else if category == "building" {
		childrenQuery := `SELECT count(id) AS childrenCount FROM rooms WHERE user_id = $1 AND parent_id = $2 AND parent_category = $3 AND deleted_at IS NULL`
		(*mockDB).ExpectQuery(regexp.QuoteMeta(childrenQuery)).
			WithArgs(testUser, testID, category).
			WillReturnRows(sqlmock.NewRows([]string{"childrenCount"}).AddRow(numberOfChildren))
	}

	if numberOfChildren == 0 {
		// Expect transaction to begin
		(*mockDB).ExpectBegin()

		// Expect the UPDATE operation
		query := fmt.Sprintf(`UPDATE "%s" SET "deleted_at"=$1 WHERE user_id = $2 AND "%s"."id" = $3 AND "%s"."deleted_at" IS NULL`, tableName, tableName, tableName)

		expectation := (*mockDB).ExpectExec(regexp.QuoteMeta(query))
		expectation.WithArgs(sqlmock.AnyArg(), testUser, testID)
		expectation.WillReturnResult(sqlmock.NewResult(0, 1))

		// Expect transaction to be committed
		(*mockDB).ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"` + testUser + `","Function":"GetAllEntities"},"Offset":"0","Limit":"15","Search":"","Filter":"[]"}`,
			`{"CacheKey":{"User":"` + testUser + `","Function":"GetAllEntities"},"Offset":"15","Limit":"15","Search":"","Filter":"[]"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"` + testUser + `","Function":"GetAllEntities"},*`).SetVal(keyVals)

		countKeys := []string{
			`{"CacheKey":{"User":"` + testUser + `","Function":"CountEntities"},"Search":"","Filter":"[]"}`,
			`{"CacheKey":{"User":"` + testUser + `","Function":"CountEntities"},"Search":"seachtext","Filter":"[]"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"` + testUser + `","Function":"CountEntities"},*`).SetVal(countKeys)

		keyVals = append(keyVals, countKeys...)

		expectedDelKeys := append(
			keyVals,
			`{"User":"`+testUser+`","Function":"GetItemParents"}`,
			`{"User":"`+testUser+`","Function":"GetContainerParents"}`,
			`{"User":"`+testUser+`","Function":"GetShelfParents"}`,
			`{"User":"`+testUser+`","Function":"GetShelving_unitParents"}`,
			`{"User":"`+testUser+`","Function":"GetRoomParents"}`,
		)

		for _, key := range expectedDelKeys {
			mockCache.ExpectDel(key).SetVal(1)
		}
	}
}

func validateDeleteEntityResponse(t *testing.T, res *http.Response, mockDB sqlmock.Sqlmock, mockCache redismock.ClientMock) {
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	contents := deleteEntitySingleResponse{}
	err = json.Unmarshal(data, &contents)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	if contents.Message != "success" {
		t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("PostGres expectations were not met: %v", err)
	}

	if err := mockCache.ExpectationsWereMet(); err != nil {
		t.Errorf("Redis expectations were not met: %v", err)
	}
}

// TestDeleteEntityValid runs the unit tests for the DeleteEntity function with valid parameters.
func TestDeleteEntityValid(t *testing.T) {
	cases := []deleteEntityTestCase{
		{
			testName:  "BEUT-109: Delete Entity Invalid Category",
			validData: false,
			id:        10,
			category:  "test",
			testUser:  "testUser1",
		},
		{
			testName:  "BEUT-110: Delete Entity Item",
			validData: true,
			id:        15,
			category:  "item",
			testUser:  "testUser2",
		},
		{
			testName:         "BEUT-111: Delete Entity Container - No Children",
			validData:        true,
			id:               20,
			category:         "container",
			testUser:         "testUser3",
			numberOfChildren: 0,
		},
		{
			testName:         "BEUT-112: Delete Entity Container - Children",
			validData:        true,
			id:               25,
			category:         "container",
			testUser:         "testUser4",
			numberOfChildren: 1,
		},
		{
			testName:         "BEUT-113: Delete Entity Shelf - No Children",
			validData:        true,
			id:               30,
			category:         "shelf",
			testUser:         "testUser5",
			numberOfChildren: 0,
		},
		{
			testName:         "BEUT-114: Delete Entity Shelf - Children",
			validData:        true,
			id:               35,
			category:         "shelf",
			testUser:         "testUser6",
			numberOfChildren: 2,
		},
		{
			testName:         "BEUT-115: Delete Entity Shelving Unit - No Children",
			validData:        true,
			id:               40,
			category:         "shelving_unit",
			testUser:         "testUser7",
			numberOfChildren: 0,
		},
		{
			testName:         "BEUT-116: Delete Entity Shelving Unit - Children",
			validData:        true,
			id:               45,
			category:         "shelving_unit",
			testUser:         "testUser8",
			numberOfChildren: 3,
		},
		{
			testName:         "BEUT-117: Delete Entity Room - No Children",
			validData:        true,
			id:               50,
			category:         "room",
			testUser:         "testUser9",
			numberOfChildren: 0,
		},
		{
			testName:         "BEUT-118: Delete Entity Room - Children",
			validData:        true,
			id:               55,
			category:         "room",
			testUser:         "testUser10",
			numberOfChildren: 4,
		},
		{
			testName:         "BEUT-119: Delete Entity Building - No Children",
			validData:        true,
			id:               60,
			category:         "building",
			testUser:         "testUser11",
			numberOfChildren: 0,
		},
		{
			testName:         "BEUT-120: Delete Entity Building - Children",
			validData:        true,
			id:               65,
			category:         "building",
			testUser:         "testUser12",
			numberOfChildren: 5,
		},
	}

	for _, tc := range cases {
		client, srv, mockDB, mockCache := setupDeleteEntityTest(t, tc.testUser)
		t.Run(tc.testName, func(t *testing.T) {

			if tc.validData {
				setupDeleteEntityMockExpectations(&mockDB, mockCache, tc.id, tc.category, tc.testUser, tc.numberOfChildren)
			}

			url := fmt.Sprintf("%s%s/%s/%d", srv.URL, deleteEntityEndpoint, tc.category, tc.id)

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

			if tc.validData && tc.numberOfChildren == 0 {
				validateDeleteEntityResponse(t, res, mockDB, mockCache)
			} else {
				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
				}
			}

		})
	}
}
