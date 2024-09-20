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
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
)

type createSingleResponse struct {
	Message string `json:"message"`
	Data    uint   `json:"ID"`
}

// TestCreateEntityItem runs our unit tests for the CreateEntity function with the shelving unit category.
func TestCreateEntityItem(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	// Create Item Test Case 1 - Valid input and includes all fields
	t.Run("CreateEntity-Item-ValidAllFields", func(t *testing.T) {
		testName := "Test Item 1"
		testNotes := "Test Notes 1"
		testCategory := "item"

		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})

	// Create Item Test Case 2 - Valid input and does not include notes
	t.Run("CreateEntity-Item-ValidNoNotes", func(t *testing.T) {
		testName := "Test Item 2"
		testCategory := "item"

		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})
}

// TestCreateEntityContainer runs our unit tests for the CreateEntity function with the shelving unit category.
func TestCreateEntityContainer(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	// Create Container Test Case 1 - Valid input and includes all fields
	t.Run("CreateEntity-Container-ValidAllFields", func(t *testing.T) {
		testName := "Test Container 1"
		testNotes := "Test Notes 1"
		testCategory := "container"

		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})

	// Create Container Test Case 2 - Valid input and does not include notes
	t.Run("CreateEntity-Container-ValidNoNotes", func(t *testing.T) {
		testName := "Test Container 2"
		testCategory := "container"

		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})
}

// TestCreateEntityShelf runs our unit tests for the CreateEntity function with the shelving unit category.
func TestCreateEntityShelf(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	// Create Shelf Test Case 1 - Valid input and includes all fields
	t.Run("CreateEntity-Shelf-ValidAllFields", func(t *testing.T) {
		testName := "Test Shelf 1"
		testNotes := "Test Notes 1"
		testCategory := "shelf"

		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})

	// Create Shelf Test Case 2 - Valid input and does not include notes
	t.Run("CreateEntity-Shelf-ValidNoNotes", func(t *testing.T) {
		testName := "Test Shelf 2"
		testCategory := "shelf"

		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})
}

// TestCreateEntityShelvingUnit runs our unit tests for the CreateEntity function with the shelving unit category.
func TestCreateEntityShelvingUnit(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	// Create ShelvingUnit Test Case 1 - Valid input and includes all fields
	t.Run("CreateEntity-ShelvingUnit-ValidAllFields", func(t *testing.T) {
		testName := "Test ShelvingUnit 1"
		testNotes := "Test Notes 1"
		testCategory := "shelvingunit"

		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})

	// Create ShelvingUnit Test Case 2 - Valid input and does not include notes
	t.Run("CreateEntity-ShelvingUnit-ValidNoNotes", func(t *testing.T) {
		testName := "Test ShelvingUnit 2"
		testCategory := "shelvingunit"

		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})
}

// TestCreateEntityRoom runs our unit tests for the CreateEntity function with the room category.
func TestCreateEntityRoom(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	// Create Room Test Case 1 - Valid input and includes all fields
	t.Run("CreateEntity-Room-ValidAllFields", func(t *testing.T) {
		testName := "Test Room 1"
		testNotes := "Test Notes 1"
		testCategory := "room"

		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})

	// Create Room Test Case 2 - Valid input and does not include notes
	t.Run("CreateEntity-Room-ValidNoNotes", func(t *testing.T) {
		testName := "Test Room 2"
		testCategory := "room"

		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})
}

// TestCreateEntityBuilding runs our unit tests for the CreateEntity function with the building category.
func TestCreateEntityBuilding(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	// Create Building Test Case 1 - Valid input and includes all fields
	t.Run("CreateEntity-Building-ValidAllFields", func(t *testing.T) {
		testName := "Test Building 1"
		testAddress := "Test Address 1"
		testNotes := "Test Notes 1"
		testCategory := "building"

		values := map[string]string{"name": testName, "address": testAddress, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})

	// Create Building Test Case 2 - Valid input and does not include notes
	t.Run("CreateEntity-Building-ValidNoNotes", func(t *testing.T) {
		testName := "Test Building 2"
		testAddress := "Test Address 2"
		testCategory := "building"

		values := map[string]string{"name": testName, "address": testAddress, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})

	// Create Building Test Case 3 - Valid input and does not include address
	t.Run("CreateEntity-Building-ValidNoAddress", func(t *testing.T) {
		testName := "Test Building 3"
		testNotes := "Test Notes 3"
		testCategory := "building"

		values := map[string]string{"name": testName, "notes": testNotes, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})

	// Create Building Test Case 4 - Valid input and does not include address or notes
	t.Run("CreateEntity-Building-ValidNoAddressNoNotes", func(t *testing.T) {
		testName := "Test Building 4"
		testCategory := "building"

		values := map[string]string{"name": testName, "category": testCategory}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateEntity(w, req)
		res := w.Result()
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

		if reflect.TypeOf(contents.Data).String() != "uint" {
			t.Errorf("Expected data to be entity id.")
		}
	})
}

// TestCreateEntityAll runs our unit tests for the CreateEntity function that apply to all categories.
func TestCreatEntityAll(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()
	// Create Entity Test Case 1 - Invalid input, missing name
	t.Run("CreateEntity-InvalidMissingName", func(t *testing.T) {

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

	// Create Entity Test Case 2 - Invalid input, missing category
	t.Run("CreateEntity-InvalidMissingCategory", func(t *testing.T) {

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

	// Create Entity Test Case 3 - Invalid input, missing name and category
	t.Run("CreateEntity-InvalidMissingNameCategory", func(t *testing.T) {

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
