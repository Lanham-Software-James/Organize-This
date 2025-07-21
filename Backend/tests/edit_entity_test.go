// Package tests is where all of out unit tests are described.
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
	"willowsuite-vault/controllers"
	"willowsuite-vault/repository"
	"willowsuite-vault/tests/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
)

type editEntitySingleResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type editEntityTestCase struct {
	testName       string `json:"-"`
	validData      bool   `json:"-"`
	EntityUser     string `json:"-"`
	EntityID       string `json:"id"`
	EntityName     string `json:"name"`
	EntityNotes    string `json:"notes"`
	EntityCategory string `json:"category"`
	EntityAddress  string `json:"address"`
	ParentID       string `json:"parentID"`
	ParentCategory string `json:"parentCategory"`
}

var editEntityEndpoint = "/v1/entity"

func setupEditEntityTest(t *testing.T, userName string) (*http.Client, *httptest.Server, sqlmock.Sqlmock, redismock.ClientMock) {
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
	r.Put("/v1/entity", handler.EditEntity)

	srv := httptest.NewServer(r)
	t.Cleanup(srv.Close)

	return &http.Client{}, srv, mockDB, mockCache
}

func setupEditEntityMockExpectations(mockDB *sqlmock.Sqlmock, mockCache redismock.ClientMock, category string, args ...string) {
	testName := args[0]
	testNotes := args[1]
	testUser := args[2]
	testID, _ := strconv.Atoi(args[3])

	tableName := category + "s"
	if category == "shelf" {
		tableName = "shelves"
	}

	// Expect transaction to begin
	(*mockDB).ExpectBegin()

	// Expect the UPDATE operation
	query := fmt.Sprintf(`UPDATE "%s" SET "name"=$1,"notes"=$2,"user_id"=$3,"created_at"=$4,"updated_at"=$5,"deleted_at"=$6,"parent_id"=$7,"parent_category"=$8 WHERE "%s"."deleted_at" IS NULL AND "id" = $9`, tableName, tableName)
	if category == "building" {
		query = `UPDATE "buildings" SET "name"=$1,"notes"=$2,"user_id"=$3,"created_at"=$4,"updated_at"=$5,"deleted_at"=$6,"address"=$7 WHERE "buildings"."deleted_at" IS NULL AND "id" = $8`
	}

	expectation := (*mockDB).ExpectExec(regexp.QuoteMeta(query))
	if category == "building" {
		testAddress := args[4]
		expectation.WithArgs(testName, testNotes, testUser, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), testAddress, testID)
	} else {
		testParentID, _ := strconv.Atoi(args[4])
		testParentCategory := args[5]
		expectation.WithArgs(testName, testNotes, testUser, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), testParentID, testParentCategory, testID)
	}
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

func validateEditEntityResponse(t *testing.T, res *http.Response, mockDB sqlmock.Sqlmock, mockCache redismock.ClientMock) {
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	contents := editEntitySingleResponse{}
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

// TestEditEntityValid runs the unit tests for the EditEntity function with valid parameters.
func TestEditEntityValid(t *testing.T) {
	cases := []editEntityTestCase{
		{
			testName:       "BEUT-63: Edit Entity Missing ID",
			validData:      false,
			EntityName:     "test item",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "container",
		},
		{
			testName:       "BEUT-64: Edit Entity Missing Name",
			validData:      false,
			EntityID:       "10",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "container",
		},
		{
			testName:       "BEUT-65: Edit Entity Missing Category",
			validData:      false,
			EntityID:       "10",
			EntityName:     "Test Container 1",
			ParentID:       "10",
			ParentCategory: "container",
		},
		{
			testName:       "BEUT-66: Edit Entity Missing Parent ID",
			validData:      false,
			EntityID:       "10",
			EntityName:     "Test Container 1",
			EntityCategory: "item",
			ParentCategory: "container",
		},
		{
			testName:       "BEUT-67: Edit Entity Missing Parent Category",
			validData:      false,
			EntityID:       "10",
			EntityName:     "Test Container 1",
			EntityCategory: "item",
			ParentID:       "10",
		},
		{
			testName:       "BEUT-68: Edit Entity Item Valid All Fields - Container",
			validData:      true,
			EntityUser:     "testuser",
			EntityID:       "10",
			EntityName:     "Test Item 1",
			EntityNotes:    "Test Notes for item 1",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "container",
		},
		{
			testName:       "BEUT-69: Edit Entity Item Valid All Fields - Shelf",
			validData:      true,
			EntityUser:     "testuser",
			EntityID:       "10",
			EntityName:     "Test Item 1",
			EntityNotes:    "Test Notes for item 1",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "shelf",
		},
		{
			testName:       "BEUT-70: Edit Entity Item Valid All Fields - Room",
			validData:      true,
			EntityUser:     "testuser",
			EntityID:       "10",
			EntityName:     "Test Item 1",
			EntityNotes:    "Test Notes for item 1",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "room",
		},
		{
			testName:       "BEUT-71: Edit Entity Item Valid No Notes",
			validData:      true,
			EntityID:       "15",
			EntityUser:     "testuser2",
			EntityName:     "Test Item 2",
			EntityNotes:    "",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "shelf",
		},
		{
			testName:       "BEUT-72: Edit Entity Item Invalid Parent Category",
			validData:      false,
			EntityName:     "Test Container 1",
			EntityCategory: "item",
			ParentID:       "10",
			ParentCategory: "test",
		},
		{
			testName:       "BEUT-73: Edit Entity Container Valid All Fields - Shelf",
			validData:      true,
			EntityID:       "20",
			EntityUser:     "testuser3",
			EntityName:     "Test Container 1",
			EntityNotes:    "Test container notes 1",
			EntityCategory: "container",
			ParentID:       "10",
			ParentCategory: "shelf",
		},
		{
			testName:       "BEUT-74: Edit Entity Container Valid All Fields - Room",
			validData:      true,
			EntityID:       "20",
			EntityUser:     "testuser3",
			EntityName:     "Test Container 1",
			EntityNotes:    "Test container notes 1",
			EntityCategory: "container",
			ParentID:       "10",
			ParentCategory: "room",
		},
		{
			testName:       "BEUT-75: Edit Entity Container Valid No Notes",
			validData:      true,
			EntityID:       "25",
			EntityUser:     "testuser4",
			EntityName:     "Test Container 2",
			EntityNotes:    "",
			EntityCategory: "container",
			ParentID:       "10",
			ParentCategory: "room",
		},
		{
			testName:       "BEUT-76: Edit Entity Container Invalid Parent Category",
			validData:      false,
			EntityID:       "25",
			EntityUser:     "testuser4",
			EntityName:     "Test Container 2",
			EntityNotes:    "",
			EntityCategory: "container",
			ParentID:       "10",
			ParentCategory: "test",
		},
		{
			testName:       "BEUT-77: Edit Entity Shelf Valid All Fields",
			validData:      true,
			EntityID:       "30",
			EntityUser:     "testuser5",
			EntityName:     "Test Shelf 1",
			EntityNotes:    "Test notes for shelf 1",
			EntityCategory: "shelf",
			ParentID:       "10",
			ParentCategory: "shelving_unit",
		},
		{
			testName:       "BEUT-78: Edit Entity Shelf Valid No Notes",
			validData:      true,
			EntityID:       "35",
			EntityUser:     "testuser6",
			EntityName:     "Test Shelf 2",
			EntityNotes:    "",
			EntityCategory: "shelf",
			ParentID:       "10",
			ParentCategory: "shelving_unit",
		},
		{
			testName:       "BEUT-79: Edit Entity Shelf Invalid Parent Category",
			validData:      false,
			EntityID:       "30",
			EntityUser:     "testuser5",
			EntityName:     "Test Shelf 1",
			EntityNotes:    "Test notes for shelf 1",
			EntityCategory: "shelf",
			ParentID:       "10",
			ParentCategory: "test",
		},
		{
			testName:       "BEUT-80: Edit Entity Shelving Unit Valid All Fields",
			validData:      true,
			EntityID:       "40",
			EntityUser:     "testuser7",
			EntityName:     "Test Shelving Unit 1",
			EntityNotes:    "Test notes for shelving unit 1",
			EntityCategory: "shelving_unit",
			ParentID:       "10",
			ParentCategory: "room",
		},
		{
			testName:       "BEUT-81: Edit Entity Shelving Unit Valid No Notes",
			validData:      true,
			EntityID:       "45",
			EntityUser:     "testuser8",
			EntityName:     "Test Shelving Unit 2",
			EntityNotes:    "",
			EntityCategory: "shelving_unit",
			ParentID:       "10",
			ParentCategory: "room",
		},
		{
			testName:       "BEUT-82: Edit Entity Shelving Unit Invalid Parent Category",
			validData:      false,
			EntityID:       "40",
			EntityUser:     "testuser7",
			EntityName:     "Test Shelving Unit 1",
			EntityNotes:    "Test notes for shelving unit 1",
			EntityCategory: "shelving_unit",
			ParentID:       "10",
			ParentCategory: "test",
		},
		{
			testName:       "BEUT-83: Edit Entity Room Valid All Fields",
			validData:      true,
			EntityID:       "50",
			EntityUser:     "testuser9",
			EntityName:     "Test Room 1",
			EntityNotes:    "Test notes for room 1",
			EntityCategory: "room",
			ParentID:       "10",
			ParentCategory: "building",
		},
		{
			testName:       "BEUT-84: Edit Entity Room Valid No Notes",
			validData:      true,
			EntityID:       "55",
			EntityUser:     "testuser10",
			EntityName:     "Test Room 2",
			EntityNotes:    "",
			EntityCategory: "room",
			ParentID:       "10",
			ParentCategory: "building",
		},
		{
			testName:       "BEUT-85: Edit Entity Room Invalid Parent Category",
			validData:      false,
			EntityID:       "50",
			EntityUser:     "testuser9",
			EntityName:     "Test Room 1",
			EntityNotes:    "Test notes for room 1",
			EntityCategory: "room",
			ParentID:       "10",
			ParentCategory: "test",
		},
		{
			testName:       "BEUT-86: Edit Entity Building Valid All Fields",
			validData:      true,
			EntityID:       "60",
			EntityUser:     "testuser11",
			EntityName:     "Test Building 1",
			EntityNotes:    "Test notes for building 1",
			EntityCategory: "building",
			EntityAddress:  "123 Test Rd",
		},
		{
			testName:       "BEUT-87: Edit Entity Building Valid No Notes",
			validData:      true,
			EntityID:       "65",
			EntityUser:     "testuser12",
			EntityName:     "Test Building 2",
			EntityNotes:    "",
			EntityCategory: "building",
			EntityAddress:  "213 Test Rd",
		},
		{
			testName:       "BEUT-88: Edit Entity Building Valid No Address",
			validData:      true,
			EntityID:       "70",
			EntityUser:     "testuser13",
			EntityName:     "Test Building 3",
			EntityNotes:    "Test notes for building 3",
			EntityCategory: "building",
			EntityAddress:  "",
		},
		{
			testName:       "BEUT-89: Edit Entity Building Valid No Notes No Address",
			validData:      true,
			EntityID:       "75",
			EntityUser:     "testuser14",
			EntityName:     "Test Building 4",
			EntityNotes:    "",
			EntityCategory: "building",
			EntityAddress:  "",
		},
	}

	for _, tc := range cases {
		client, srv, mockDB, mockCache := setupEditEntityTest(t, tc.EntityUser)
		t.Run(tc.testName, func(t *testing.T) {
			if tc.EntityCategory == "building" && tc.validData {
				setupEditEntityMockExpectations(&mockDB, mockCache, tc.EntityCategory, tc.EntityName, tc.EntityNotes, tc.EntityUser, tc.EntityID, tc.EntityAddress)
			} else if tc.validData {
				setupEditEntityMockExpectations(&mockDB, mockCache, tc.EntityCategory, tc.EntityName, tc.EntityNotes, tc.EntityUser, tc.EntityID, tc.ParentID, tc.ParentCategory)
			}

			url := fmt.Sprintf("%s%s", srv.URL, editEntityEndpoint)

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

			if tc.validData {
				validateEditEntityResponse(t, res, mockDB, mockCache)
			} else {
				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
				}
			}

		})
	}
}
