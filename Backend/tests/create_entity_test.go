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
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
)

type createSingleResponse struct {
	Message string `json:"message"`
	Data    uint   `json:"ID"`
}

type createEntityTestCase struct {
	testName       string
	entityName     string
	entityNotes    string
	entityCategory string
	entityUser     string
	entityID       uint
	entityAddress  string
}

var endpoint = "/v1/entity"

func setupCreateEntityTest(t *testing.T, userName string) (*http.Client, *httptest.Server, sqlmock.Sqlmock, redismock.ClientMock) {
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
	r.Post("/v1/entity", handler.CreateEntity)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &http.Client{}, srv, mockDB, mockCache
}

func setupCreateEntityMockExpectations(mockDB *sqlmock.Sqlmock, mockCache redismock.ClientMock, category string, args ...string) {
	(*mockDB).ExpectBegin()

	testName := args[0]
	testNotes := args[1]
	testUser := args[2]
	testID := args[3]

	tableName := category + "s"
	if category == "shelvingunit" {
		tableName = "shelving_units"
	} else if category == "shelf" {
		tableName = "shelves"
	}

	query := fmt.Sprintf(`INSERT INTO "%s" \("name","notes","user_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`, tableName)
	if category == "building" {
		query = `INSERT INTO "buildings" \("name","notes","user_id","created_at","updated_at","address"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`
	}

	expectation := (*mockDB).ExpectQuery(query)
	if category == "building" {
		testAddress := args[4]
		expectation.WithArgs(testName, testNotes, testUser, sqlmock.AnyArg(), sqlmock.AnyArg(), testAddress)
	} else {
		expectation.WithArgs(testName, testNotes, testUser, sqlmock.AnyArg(), sqlmock.AnyArg())
	}
	expectation.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

	(*mockDB).ExpectCommit()

	keyVals := []string{
		`{"CacheKey":{"User":"` + testUser + `","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
		`{"CacheKey":{"User":"` + testUser + `","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
	}
	mockCache.ExpectKeys(`{"CacheKey":{"User":"` + testUser + `","Function":"GetAllEntities"},*`).SetVal(keyVals)
	mockCache.ExpectDel(`{"CacheKey":{"User":"` + testUser + `","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
	mockCache.ExpectDel(`{"CacheKey":{"User":"` + testUser + `","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
	mockCache.ExpectDel(`{"User":"` + testUser + `","Function":"CountEntities"}`).SetVal(1)
}

func validateCreateEntitySuccessResponse(t *testing.T, res *http.Response, mockDB sqlmock.Sqlmock, mockCache redismock.ClientMock, testID uint) {
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	contents := createSingleResponse{}
	err = json.Unmarshal(data, &contents)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	if contents.Message != "success" {
		t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
	}

	if contents.Data == testID {
		t.Errorf("Expected data to be %v. Got: %v", testID, contents.Data)
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("PostGres expectations were not met: %v", err)
	}

	if err := mockCache.ExpectationsWereMet(); err != nil {
		t.Errorf("Redis expectations were not met: %v", err)
	}
}

// TestCreateEntityInvalid runs the unit tests for invalid cases.
func TestCreateEntityInvalid(t *testing.T) {
	cases := []createEntityTestCase{
		{
			testName:       "BEUT-1: Create Entity Missing Name",
			entityCategory: "item",
		},
		{
			testName:   "BEUT-2: Create Entity Missing Category",
			entityName: "Test Container 1",
		},
		{
			testName: "BEUT-3: Create Entity Missing Name and Category",
		},
	}

	for _, tc := range cases {
		client, srv, mockDB, mockCache := setupCreateEntityTest(t, tc.entityUser)
		t.Run(tc.testName, func(t *testing.T) {
			if tc.entityCategory == "building" {
				setupCreateEntityMockExpectations(&mockDB, mockCache, tc.entityCategory, tc.entityName, tc.entityNotes, tc.entityUser, strconv.Itoa(int(tc.entityID)), tc.entityAddress)
			} else {
				setupCreateEntityMockExpectations(&mockDB, mockCache, tc.entityCategory, tc.entityName, tc.entityNotes, tc.entityUser, strconv.Itoa(int(tc.entityID)))
			}

			values := map[string]string{"name": tc.entityName, "notes": tc.entityNotes, "category": tc.entityCategory}
			if tc.entityCategory == "building" {
				values = map[string]string{"name": tc.entityName, "notes": tc.entityNotes, "category": tc.entityCategory, "address": tc.entityAddress}
			}

			payload, err := json.Marshal(values)
			if err != nil {
				t.Fatalf("Failed to marshal json: %v", err)
			}

			res, err := client.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusBadRequest {
				t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
			}
		})
	}
}

// TestCreateEntityValid runs the unit tests for the CreateEntity function with valid parameters.
func TestCreateEntityValid(t *testing.T) {
	cases := []createEntityTestCase{
		{
			testName:       "BEUT-4: Create Entity Item Valid All Fields",
			entityName:     "Test Item 1",
			entityNotes:    "Test Notes for item 1",
			entityCategory: "item",
			entityUser:     "testuser",
			entityID:       10,
		},
		{
			testName:       "BEUT-5: Create Entity Item Valid No Notes",
			entityName:     "Test Item 2",
			entityNotes:    "",
			entityCategory: "item",
			entityUser:     "testuser2",
			entityID:       15,
		},
		{
			testName:       "BEUT-6: Create Entity Container Valid All Fields",
			entityName:     "Test Container 1",
			entityNotes:    "Test container notes 1",
			entityCategory: "container",
			entityUser:     "testuser3",
			entityID:       20,
		},
		{
			testName:       "BEUT-7: Create Entity Container Valid No Notes",
			entityName:     "Test Container 2",
			entityNotes:    "",
			entityCategory: "container",
			entityUser:     "testuser4",
			entityID:       25,
		},
		{
			testName:       "BEUT-8: Create Entity Shelf Valid All Fields",
			entityName:     "Test Shelf 1",
			entityNotes:    "Test notes for shelf 1",
			entityCategory: "shelf",
			entityUser:     "testuser5",
			entityID:       30,
		},
		{
			testName:       "BEUT-9: Create Entity Shelf Valid No Notes",
			entityName:     "Test Shelf 2",
			entityNotes:    "",
			entityCategory: "shelf",
			entityUser:     "testuser6",
			entityID:       35,
		},
		{
			testName:       "BEUT-10: Create Entity Shelving Unit Valid All Fields",
			entityName:     "Test Shelving Unit 1",
			entityNotes:    "Test notes for shelving unit 1",
			entityCategory: "shelvingunit",
			entityUser:     "testuser7",
			entityID:       40,
		},
		{
			testName:       "BEUT-11: Create Entity Shelving Unit Valid No Notes",
			entityName:     "Test Shelving Unit 2",
			entityNotes:    "",
			entityCategory: "shelvingunit",
			entityUser:     "testuser8",
			entityID:       45,
		},
		{
			testName:       "BEUT-12: Create Entity Room Valid All Fields",
			entityName:     "Test Room 1",
			entityNotes:    "Test notes for room 1",
			entityCategory: "room",
			entityUser:     "testuser9",
			entityID:       50,
		},
		{
			testName:       "BEUT-13: Create Entity Room Valid No Notes",
			entityName:     "Test Room 2",
			entityNotes:    "",
			entityCategory: "room",
			entityUser:     "testuser10",
			entityID:       55,
		},
		{
			testName:       "BEUT-14: Create Entity Building Valid All Fields",
			entityName:     "Test Building 1",
			entityNotes:    "Test notes for building 1",
			entityCategory: "building",
			entityUser:     "testuser11",
			entityID:       60,
			entityAddress:  "123 Test Rd",
		},
		{
			testName:       "BEUT-15: Create Entity Building Valid No Notes",
			entityName:     "Test Building 2",
			entityNotes:    "",
			entityCategory: "building",
			entityUser:     "testuser12",
			entityID:       65,
			entityAddress:  "213 Test Rd",
		},
		{
			testName:       "BEUT-16: Create Entity Building Valid No Address",
			entityName:     "Test Building 3",
			entityNotes:    "Test notes for building 3",
			entityCategory: "building",
			entityUser:     "testuser13",
			entityID:       70,
			entityAddress:  "",
		},
		{
			testName:       "BEUT-17: Create Entity Building Valid No Notes No Address",
			entityName:     "Test Building 4",
			entityNotes:    "",
			entityCategory: "building",
			entityUser:     "testuser14",
			entityID:       75,
			entityAddress:  "",
		},
	}

	for _, tc := range cases {
		client, srv, mockDB, mockCache := setupCreateEntityTest(t, tc.entityUser)
		t.Run(tc.testName, func(t *testing.T) {
			if tc.entityCategory == "building" {
				setupCreateEntityMockExpectations(&mockDB, mockCache, tc.entityCategory, tc.entityName, tc.entityNotes, tc.entityUser, strconv.Itoa(int(tc.entityID)), tc.entityAddress)
			} else {
				setupCreateEntityMockExpectations(&mockDB, mockCache, tc.entityCategory, tc.entityName, tc.entityNotes, tc.entityUser, strconv.Itoa(int(tc.entityID)))
			}

			values := map[string]string{"name": tc.entityName, "notes": tc.entityNotes, "category": tc.entityCategory}
			if tc.entityCategory == "building" {
				values = map[string]string{"name": tc.entityName, "notes": tc.entityNotes, "category": tc.entityCategory, "address": tc.entityAddress}
			}

			payload, err := json.Marshal(values)
			if err != nil {
				t.Fatalf("Failed to marshal json: %v", err)
			}

			res, err := client.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer res.Body.Close()

			validateCreateEntitySuccessResponse(t, res, mockDB, mockCache, tc.entityID)
		})
	}
}
