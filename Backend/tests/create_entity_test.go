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
	testName       string `json:"-"`
	validData      bool   `json:"-"`
	EntityUser     string `json:"-"`
	EntityID       uint   `json:"-"`
	EntityName     string `json:"name"`
	EntityNotes    string `json:"notes"`
	EntityCategory string `json:"category"`
	EntityAddress  string `json:"address"`
	ParentID       string `json:"parentID"`
	ParentCategory string `json:"parentCategory"`
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
	if category == "shelf" {
		tableName = "shelves"
	}

	query := fmt.Sprintf(`INSERT INTO "%s" ("name","notes","user_id","created_at","updated_at","deleted_at","parent_id","parent_category") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`, tableName)
	if category == "building" {
		query = `INSERT INTO "buildings" ("name","notes","user_id","created_at","updated_at","deleted_at","address") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`
	}

	expectation := (*mockDB).ExpectQuery(regexp.QuoteMeta(query))
	if category == "building" {
		testAddress := args[4]
		expectation.WithArgs(testName, testNotes, testUser, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), testAddress)
	} else {
		testParentID, _ := strconv.Atoi(args[4])
		testParentCategory := args[5]
		expectation.WithArgs(testName, testNotes, testUser, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), testParentID, testParentCategory)
	}
	expectation.WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

	(*mockDB).ExpectCommit()

	keyVals := []string{
		`{"CacheKey":{"User":"` + testUser + `","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
		`{"CacheKey":{"User":"` + testUser + `","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
	}

	mockCache.ExpectKeys(`{"CacheKey":{"User":"` + testUser + `","Function":"GetAllEntities"},*`).SetVal(keyVals)

	expectedDelKeys := append(
		keyVals,
		`{"User":"`+testUser+`","Function":"CountEntities"}`,
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

// TestCreateEntityValid runs the unit tests for the CreateEntity function with valid parameters.
func TestCreateEntityValid(t *testing.T) {
	cases := []createEntityTestCase{
		{
			testName:       "BEUT-1: Create Entity Missing Name",
			validData:      false,
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "container",
		},
		{
			testName:       "BEUT-2: Create Entity Missing Category",
			validData:      false,
			EntityName:     "Test Container 1",
			ParentID:       "10",
			ParentCategory: "container",
		},
		{
			testName:       "BEUT-3: Create Entity Missing Parent ID",
			validData:      false,
			EntityName:     "Test Container 1",
			EntityCategory: "item",
			ParentCategory: "container",
		},
		{
			testName:       "BEUT-4: Create Entity Missing Parent Category",
			validData:      false,
			EntityName:     "Test Container 1",
			EntityCategory: "item",
			ParentID:       "10",
		},
		{
			testName:       "BEUT-5: Create Entity Item Valid All Fields - Container",
			validData:      true,
			EntityUser:     "testuser",
			EntityID:       10,
			EntityName:     "Test Item 1",
			EntityNotes:    "Test Notes for item 1",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "container",
		},
		{
			testName:       "BEUT-6: Create Entity Item Valid All Fields - Shelf",
			validData:      true,
			EntityUser:     "testuser",
			EntityID:       10,
			EntityName:     "Test Item 1",
			EntityNotes:    "Test Notes for item 1",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "shelf",
		},
		{
			testName:       "BEUT-7: Create Entity Item Valid All Fields - Room",
			validData:      true,
			EntityUser:     "testuser",
			EntityID:       10,
			EntityName:     "Test Item 1",
			EntityNotes:    "Test Notes for item 1",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "room",
		},
		{
			testName:       "BEUT-8: Create Entity Item Valid No Notes",
			validData:      true,
			EntityID:       15,
			EntityUser:     "testuser2",
			EntityName:     "Test Item 2",
			EntityNotes:    "",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "shelf",
		},
		{
			testName:       "BEUT-9: Create Entity Item Invalid Parent Category",
			validData:      false,
			EntityName:     "Test Container 1",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "test",
		},
		{
			testName:       "BEUT-10: Create Entity Container Valid All Fields - Shelf",
			validData:      true,
			EntityID:       20,
			EntityUser:     "testuser3",
			EntityName:     "Test Container 1",
			EntityNotes:    "Test container notes 1",
			EntityCategory: "container",
			ParentID:       "10",
			ParentCategory: "shelf",
		},
		{
			testName:       "BEUT-11: Create Entity Container Valid All Fields - Room",
			validData:      true,
			EntityID:       20,
			EntityUser:     "testuser3",
			EntityName:     "Test Container 1",
			EntityNotes:    "Test container notes 1",
			EntityCategory: "container",
			ParentID:       "10",
			ParentCategory: "room",
		},
		{
			testName:       "BEUT-12: Create Entity Container Valid No Notes",
			validData:      true,
			EntityID:       25,
			EntityUser:     "testuser4",
			EntityName:     "Test Container 2",
			EntityNotes:    "",
			EntityCategory: "container",
			ParentID:       "10",
			ParentCategory: "room",
		},
		{
			testName:       "BEUT-13: Create Entity Container Invalid Parent Category",
			validData:      false,
			EntityID:       25,
			EntityUser:     "testuser4",
			EntityName:     "Test Container 2",
			EntityNotes:    "",
			EntityCategory: "container",
			ParentID:       "10",
			ParentCategory: "test",
		},
		{
			testName:       "BEUT-14: Create Entity Shelf Valid All Fields",
			validData:      true,
			EntityID:       30,
			EntityUser:     "testuser5",
			EntityName:     "Test Shelf 1",
			EntityNotes:    "Test notes for shelf 1",
			EntityCategory: "shelf",
			ParentID:       "10",
			ParentCategory: "shelving_unit",
		},
		{
			testName:       "BEUT-15: Create Entity Shelf Valid No Notes",
			validData:      true,
			EntityID:       35,
			EntityUser:     "testuser6",
			EntityName:     "Test Shelf 2",
			EntityNotes:    "",
			EntityCategory: "shelf",
			ParentID:       "10",
			ParentCategory: "shelving_unit",
		},
		{
			testName:       "BEUT-16: Create Entity Shelf Invalid Parent Category",
			validData:      false,
			EntityID:       30,
			EntityUser:     "testuser5",
			EntityName:     "Test Shelf 1",
			EntityNotes:    "Test notes for shelf 1",
			EntityCategory: "shelf",
			ParentID:       "10",
			ParentCategory: "test",
		},
		{
			testName:       "BEUT-17: Create Entity Shelving Unit Valid All Fields",
			validData:      true,
			EntityID:       40,
			EntityUser:     "testuser7",
			EntityName:     "Test Shelving Unit 1",
			EntityNotes:    "Test notes for shelving unit 1",
			EntityCategory: "shelving_unit",
			ParentID:       "10",
			ParentCategory: "room",
		},
		{
			testName:       "BEUT-18: Create Entity Shelving Unit Valid No Notes",
			validData:      true,
			EntityID:       45,
			EntityUser:     "testuser8",
			EntityName:     "Test Shelving Unit 2",
			EntityNotes:    "",
			EntityCategory: "shelving_unit",
			ParentID:       "10",
			ParentCategory: "room",
		},
		{
			testName:       "BEUT-19: Create Entity Shelving Unit Invalid Parent Category",
			validData:      false,
			EntityID:       40,
			EntityUser:     "testuser7",
			EntityName:     "Test Shelving Unit 1",
			EntityNotes:    "Test notes for shelving unit 1",
			EntityCategory: "shelving_unit",
			ParentID:       "10",
			ParentCategory: "test",
		},
		{
			testName:       "BEUT-20: Create Entity Room Valid All Fields",
			validData:      true,
			EntityID:       50,
			EntityUser:     "testuser9",
			EntityName:     "Test Room 1",
			EntityNotes:    "Test notes for room 1",
			EntityCategory: "room",
			ParentID:       "10",
			ParentCategory: "building",
		},
		{
			testName:       "BEUT-21: Create Entity Room Valid No Notes",
			validData:      true,
			EntityID:       55,
			EntityUser:     "testuser10",
			EntityName:     "Test Room 2",
			EntityNotes:    "",
			EntityCategory: "room",
			ParentID:       "10",
			ParentCategory: "building",
		},
		{
			testName:       "BEUT-22: Create Entity Room Invalid Parent Category",
			validData:      false,
			EntityID:       50,
			EntityUser:     "testuser9",
			EntityName:     "Test Room 1",
			EntityNotes:    "Test notes for room 1",
			EntityCategory: "room",
			ParentID:       "10",
			ParentCategory: "test",
		},
		{
			testName:       "BEUT-23: Create Entity Building Valid All Fields",
			validData:      true,
			EntityID:       60,
			EntityUser:     "testuser11",
			EntityName:     "Test Building 1",
			EntityNotes:    "Test notes for building 1",
			EntityCategory: "building",
			EntityAddress:  "123 Test Rd",
		},
		{
			testName:       "BEUT-24: Create Entity Building Valid No Notes",
			validData:      true,
			EntityID:       65,
			EntityUser:     "testuser12",
			EntityName:     "Test Building 2",
			EntityNotes:    "",
			EntityCategory: "building",
			EntityAddress:  "213 Test Rd",
		},
		{
			testName:       "BEUT-25: Create Entity Building Valid No Address",
			validData:      true,
			EntityID:       70,
			EntityUser:     "testuser13",
			EntityName:     "Test Building 3",
			EntityNotes:    "Test notes for building 3",
			EntityCategory: "building",
			EntityAddress:  "",
		},
		{
			testName:       "BEUT-26: Create Entity Building Valid No Notes No Address",
			validData:      true,
			EntityID:       75,
			EntityUser:     "testuser14",
			EntityName:     "Test Building 4",
			EntityNotes:    "",
			EntityCategory: "building",
			EntityAddress:  "",
		},
	}

	for _, tc := range cases {
		client, srv, mockDB, mockCache := setupCreateEntityTest(t, tc.EntityUser)
		t.Run(tc.testName, func(t *testing.T) {
			if tc.EntityCategory == "building" && tc.validData {
				setupCreateEntityMockExpectations(&mockDB, mockCache, tc.EntityCategory, tc.EntityName, tc.EntityNotes, tc.EntityUser, strconv.Itoa(int(tc.EntityID)), tc.EntityAddress)
			} else if tc.validData {
				setupCreateEntityMockExpectations(&mockDB, mockCache, tc.EntityCategory, tc.EntityName, tc.EntityNotes, tc.EntityUser, strconv.Itoa(int(tc.EntityID)), tc.ParentID, tc.ParentCategory)
			}

			payload, err := json.Marshal(tc)
			if err != nil {
				t.Fatalf("Failed to marshal json: %v", err)
			}

			res, err := client.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer res.Body.Close()

			if tc.validData {
				validateCreateEntitySuccessResponse(t, res, mockDB, mockCache, tc.EntityID)
			} else {
				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
				}
			}

		})
	}
}
