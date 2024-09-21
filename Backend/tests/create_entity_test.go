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

// TestCreateEntityAll runs our unit tests for the CreateEntity function that apply to all categories.
func TestCreatEntityAll(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
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
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateEntity)

	srv := httptest.NewServer(r)
	defer srv.Close()

	t.Run("BEUT-4: Create Entity Item Valid All Fields", func(t *testing.T) {
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

	t.Run("BEUT-5: Create Entity Item Valid No Notes", func(t *testing.T) {
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

	t.Run("BEUT-6: Create Entity Container Valid All Fields", func(t *testing.T) {
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

	t.Run("BEUT-7: Create Entity Container Valid No Notes", func(t *testing.T) {
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

	t.Run("BEUT-8: Create Entity Shelf Valid All Fields", func(t *testing.T) {
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

	t.Run("BEUT-9: Create Entity Shelf Valid No Notes", func(t *testing.T) {
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

	t.Run("BEUT-10: Create Entity Shelving Unit Valid All Fields", func(t *testing.T) {
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

	t.Run("BEUT-11: Create Entity Shelving Unit Valid No Notes", func(t *testing.T) {
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

	t.Run("BEUT-12: Create Entity Room Valid All Fields", func(t *testing.T) {
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

	t.Run("BEUT-13: Create Entity Room Valid No Notes", func(t *testing.T) {
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

	t.Run("BEUT-14: Create Entity Building Valid All Fields", func(t *testing.T) {
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

	t.Run("BEUT-15: Create Entity Building Valid No Notes", func(t *testing.T) {
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

	t.Run("BEUT-16: Create Entity Building Valid No Address", func(t *testing.T) {
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

	t.Run("BEUT-17: Create Entity Building Valid Not Notes No Address ", func(t *testing.T) {
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
