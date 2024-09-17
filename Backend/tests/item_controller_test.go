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
	"testing"

	"github.com/go-chi/chi/v5"
)

type itemSingleResponse struct {
	Message string `json:"message"`
	Data    struct {
		ID        uint   `json:"ID"`
		Name      string `json:"Name"`
		Notes     string `json:"Notes"`
		CreatedAt string `json:"CreatedAt"`
		UpdatedAt string `json:"UpdatedAt"`
	}
}

// TestCreatItem runs all our various unit tests for the CreateItem function
func TestCreateItem(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity-management/item"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateItem)

	srv := httptest.NewServer(r)
	defer srv.Close()

	// Create Item Test Case 1 - Valid input and includes all fields
	t.Run("CreateItem-ValidAllFields", func(t *testing.T) {
		testName := "Test Item 1"
		testNotes := "Test Notes 1"

		values := map[string]string{"Name": testName, "Notes": testNotes}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateItem(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		contents := itemSingleResponse{}
		err = json.Unmarshal(data, &contents)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		if contents.Message != "success" {
			t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
		}

		if contents.Data.Name != testName {
			t.Errorf("Expected item name to be: %s. Got: %s.", testName, contents.Data.Name)
		}

		if contents.Data.Notes != testNotes {
			t.Errorf("Expected item notes to be: %s. Got: %s", testNotes, contents.Data.Notes)
		}
	})

	// Create Item Test Case 2 - Valid input and does not include notes
	t.Run("CreateItem-ValidNoNotes", func(t *testing.T) {
		testName := "Test Item 2"

		values := map[string]string{"Name": testName}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateItem(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		contents := itemSingleResponse{}
		err = json.Unmarshal(data, &contents)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		if contents.Message != "success" {
			t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
		}

		if contents.Data.Name != testName {
			t.Errorf("Expected item name to be: %s. Got: %s.", testName, contents.Data.Name)
		}

		if contents.Data.Notes != "" {
			t.Errorf("Expected item notes to be: null. Got: %s", contents.Data.Notes)
		}
	})

	// Create Item Test Case 3 - Invalid input, missing name
	t.Run("CreateItem-InvalidMissingName", func(t *testing.T) {

		values := map[string]string{}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateItem(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})
}
