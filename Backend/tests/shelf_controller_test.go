package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"organize-this/controllers"
	"organize-this/repository"
	"testing"

	"github.com/go-chi/chi/v5"
)

type shelfSingleResponse struct {
	Message string `json:"message"`
	Data    struct {
		ID        uint   `json:"ID"`
		Name      string `json:"Name"`
		Notes     string `json:"Notes"`
		CreatedAt string `json:"CreatedAt"`
		UpdatedAt string `json:"UpdatedAt"`
	}
}

func TestCreateShelf(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity-management/shelf"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateShelf)

	srv := httptest.NewServer(r)
	defer srv.Close()

	// Create Shelf Test Case 1 - Valid input and includes all fields
	t.Run("CreateShelf-ValidAllFields", func(t *testing.T) {
		testName := "Test Shelf 1"
		testNotes := "Test Notes 1"

		values := map[string]string{"Name": testName, "Notes": testNotes}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateShelf(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		contents := shelfSingleResponse{}
		err = json.Unmarshal(data, &contents)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		if contents.Message != "success" {
			t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
		}

		if contents.Data.Name != testName {
			t.Errorf("Expected shelf name to be: %s. Got: %s.", testName, contents.Data.Name)
		}

		if contents.Data.Notes != testNotes {
			t.Errorf("Expected shelf notes to be: %s. Got: %s", testNotes, contents.Data.Notes)
		}
	})

	// Create Shelf Test Case 2 - Valid input and does not include notes
	t.Run("CreateShelf-ValidNoNotes", func(t *testing.T) {
		testName := "Test Shelf 2"

		values := map[string]string{"Name": testName}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateShelf(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		contents := shelfSingleResponse{}
		err = json.Unmarshal(data, &contents)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		if contents.Message != "success" {
			t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
		}

		if contents.Data.Name != testName {
			t.Errorf("Expected shelf name to be: %s. Got: %s.", testName, contents.Data.Name)
		}

		if contents.Data.Notes != "" {
			t.Errorf("Expected shelf notes to be: null. Got: %s", contents.Data.Notes)
		}
	})

	// Create Shelf Test Case 3 - Invalid input, missing name
	t.Run("CreateShelf-InvalidMissingName", func(t *testing.T) {

		values := map[string]string{}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateShelf(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})
}
