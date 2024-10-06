// Package tests is where all of out unit tests are described.
package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"organize-this/controllers"
	"organize-this/repository"
	"organize-this/tests/mocks"
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

// TestCreateEntityAll runs our unit tests for the CreateEntity function that apply to all categories.
func TestCreatEntityAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, _ := mocks.NewMockDB()
	mockCache, _ := redismock.NewClientMock()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB, Cache: mockCache}, CognitoClient: mocks.NewMockCognitoClient(ctrl)}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	t.Run("BEUT-1: Create Entity Missing Name", func(t *testing.T) {

		testCategory := "building"
		values := map[string]string{"category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("BEUT-2: Create Entity Missing Category", func(t *testing.T) {

		testName := "Test Building 6"
		values := map[string]string{"name": testName}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("BEUT-3: Create Entity Missing Name and Category", func(t *testing.T) {

		values := map[string]string{}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})
}

// TestCreateEntityItem runs our unit tests for the CreateEntity function with the shelving unit category.
func TestCreateEntityItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postgres, mockDB := mocks.NewMockDB()
	redis, mockCache := redismock.NewClientMock()
	handler := controllers.Handler{Repository: &repository.Repository{Database: postgres, Cache: redis}, CognitoClient: mocks.NewMockCognitoClient(ctrl)}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Use(mocks.MockJWTMiddleware)
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	t.Run("BEUT-4: Create Entity Item Valid All Fields", func(t *testing.T) {
		testName := "Test Item 1"
		testNotes := "Test Notes 1"
		testCategory := "item"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "items" \("name","notes","user_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
			WithArgs(testName, testNotes, "testuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))
		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)

		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})

	t.Run("BEUT-5: Create Entity Item Valid No Notes", func(t *testing.T) {
		testName := "Test Item 2"
		testCategory := "item"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "items" \("name","notes","user_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
			WithArgs(testName, "", "testuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)

		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})
}

// TestCreateEntityContainer runs our unit tests for the CreateEntity function with the shelving unit category.
func TestCreateEntityContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postgres, mockDB := mocks.NewMockDB()
	redis, mockCache := redismock.NewClientMock()
	handler := controllers.Handler{Repository: &repository.Repository{Database: postgres, Cache: redis}, CognitoClient: mocks.NewMockCognitoClient(ctrl)}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Use(mocks.MockJWTMiddleware)
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	t.Run("BEUT-6: Create Entity Container Valid All Fields", func(t *testing.T) {
		testName := "Test Container 1"
		testNotes := "Test Notes 1"
		testCategory := "container"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "containers" \("name","notes","user_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
			WithArgs(testName, testNotes, "testuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)

		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})

	t.Run("BEUT-7: Create Entity Container Valid No Notes", func(t *testing.T) {
		testName := "Test Container 2"
		testCategory := "container"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "containers" \("name","notes","user_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
			WithArgs(testName, "", "testuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)

		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})
}

// TestCreateEntityShelf runs our unit tests for the CreateEntity function with the shelving unit category.
func TestCreateEntityShelf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postgres, mockDB := mocks.NewMockDB()
	redis, mockCache := redismock.NewClientMock()
	handler := controllers.Handler{Repository: &repository.Repository{Database: postgres, Cache: redis}, CognitoClient: mocks.NewMockCognitoClient(ctrl)}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Use(mocks.MockJWTMiddleware)
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	t.Run("BEUT-8: Create Entity Shelf Valid All Fields", func(t *testing.T) {
		testName := "Test Shelf 1"
		testNotes := "Test Notes 1"
		testCategory := "shelf"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "shelves" \("name","notes","user_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
			WithArgs(testName, testNotes, "testuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)

		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})

	t.Run("BEUT-9: Create Entity Shelf Valid No Notes", func(t *testing.T) {
		testName := "Test Shelf 2"
		testCategory := "shelf"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "shelves" \("name","notes","user_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
			WithArgs(testName, "", "testuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)

		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})
}

// TestCreateEntityShelvingUnit runs our unit tests for the CreateEntity function with the shelving unit category.
func TestCreateEntityShelvingUnit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postgres, mockDB := mocks.NewMockDB()
	redis, mockCache := redismock.NewClientMock()
	handler := controllers.Handler{Repository: &repository.Repository{Database: postgres, Cache: redis}, CognitoClient: mocks.NewMockCognitoClient(ctrl)}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Use(mocks.MockJWTMiddleware)
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	t.Run("BEUT-10: Create Entity Shelving Unit Valid All Fields", func(t *testing.T) {
		testName := "Test ShelvingUnit 1"
		testNotes := "Test Notes 1"
		testCategory := "shelvingunit"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "shelving_units" \("name","notes","user_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
			WithArgs(testName, testNotes, "testuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)

		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})

	t.Run("BEUT-11: Create Entity Shelving Unit Valid No Notes", func(t *testing.T) {
		testName := "Test ShelvingUnit 2"
		testCategory := "shelvingunit"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "shelving_units" \("name","notes","user_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
			WithArgs(testName, "", "testuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)

		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})
}

// TestCreateEntityRoom runs our unit tests for the CreateEntity function with the room category.
func TestCreateEntityRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postgres, mockDB := mocks.NewMockDB()
	redis, mockCache := redismock.NewClientMock()
	handler := controllers.Handler{Repository: &repository.Repository{Database: postgres, Cache: redis}, CognitoClient: mocks.NewMockCognitoClient(ctrl)}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Use(mocks.MockJWTMiddleware)
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	t.Run("BEUT-12: Create Entity Room Valid All Fields", func(t *testing.T) {
		testName := "Test Room 1"
		testNotes := "Test Notes 1"
		testCategory := "room"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "rooms" \("name","notes","user_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
			WithArgs(testName, testNotes, "testuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)

		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})

	t.Run("BEUT-13: Create Entity Room Valid No Notes", func(t *testing.T) {
		testName := "Test Room 2"
		testCategory := "room"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "rooms" \("name","notes","user_id","created_at","updated_at"\) VALUES \(\$1,\$2,\$3,\$4,\$5\) RETURNING "id"`).
			WithArgs(testName, "", "testuser", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)

		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})
}

// TestCreateEntityBuilding runs our unit tests for the CreateEntity function with the building category.
func TestCreateEntityBuilding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postgres, mockDB := mocks.NewMockDB()
	redis, mockCache := redismock.NewClientMock()
	handler := controllers.Handler{Repository: &repository.Repository{Database: postgres, Cache: redis}, CognitoClient: mocks.NewMockCognitoClient(ctrl)}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Use(mocks.MockJWTMiddleware)
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	t.Run("BEUT-14: Create Entity Building Valid All Fields", func(t *testing.T) {
		testName := "Test Building 1"
		testAddress := "Test Address 1"
		testNotes := "Test Notes 1"
		testCategory := "building"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "buildings" \("name","notes","user_id","created_at","updated_at","address"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
			WithArgs(testName, testNotes, "testuser", sqlmock.AnyArg(), sqlmock.AnyArg(), testAddress).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)
		values := map[string]string{"name": testName, "address": testAddress, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})

	t.Run("BEUT-15: Create Entity Building Valid No Notes", func(t *testing.T) {
		testName := "Test Building 2"
		testAddress := "Test Address 2"
		testCategory := "building"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "buildings" \("name","notes","user_id","created_at","updated_at","address"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
			WithArgs(testName, "", "testuser", sqlmock.AnyArg(), sqlmock.AnyArg(), testAddress).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)
		values := map[string]string{"name": testName, "address": testAddress, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})

	t.Run("BEUT-16: Create Entity Building Valid No Address", func(t *testing.T) {
		testName := "Test Building 3"
		testNotes := "Test Notes 3"
		testCategory := "building"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "buildings" \("name","notes","user_id","created_at","updated_at","address"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
			WithArgs(testName, testNotes, "testuser", sqlmock.AnyArg(), sqlmock.AnyArg(), "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)
		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})

	t.Run("BEUT-17: Create Entity Building Valid Not Notes No Address ", func(t *testing.T) {
		testName := "Test Building 4"
		testCategory := "building"
		var testID uint = 1

		mockDB.ExpectBegin()
		mockDB.ExpectQuery(`INSERT INTO "buildings" \("name","notes","user_id","created_at","updated_at","address"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\) RETURNING "id"`).
			WithArgs(testName, "", "testuser", sqlmock.AnyArg(), sqlmock.AnyArg(), "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testID))

		mockDB.ExpectCommit()

		keyVals := []string{
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`,
			`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`,
		}
		mockCache.ExpectKeys(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},*`).SetVal(keyVals)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"0","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"CacheKey":{"User":"testuser","Function":"GetAllEntities"},"Offset":"15","Limit":"15"}`).SetVal(1)
		mockCache.ExpectDel(`{"User":"testuser","Function":"CountEntities"}`).SetVal(1)
		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		res, err := http.Post(srv.URL+endpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer res.Body.Close()

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
	})
}
