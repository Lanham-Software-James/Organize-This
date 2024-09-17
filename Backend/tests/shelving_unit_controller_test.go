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

type shelvingUnitSingleResponse struct {
	Message string `json:"message"`
	Data    struct {
		ID        uint   `json:"ID"`
		Name      string `json:"Name"`
		Notes     string `json:"Notes"`
		CreatedAt string `json:"CreatedAt"`
		UpdatedAt string `json:"UpdatedAt"`
	}
}

func TestCreateShelvingUnit(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity-management/shelvingunit"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateShelvingUnit)

	srv := httptest.NewServer(r)
	defer srv.Close()

	// Create Shelving Unit Test Case 1 - Valid input and includes all fields
	t.Run("CreateShelvingUnit-ValidAllFields", func(t *testing.T) {
		testName := "Test Shelving Unit 1"
		testNotes := "Test Notes 1"

		values := map[string]string{"Name": testName, "Notes": testNotes}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateShelvingUnit(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		contents := shelvingUnitSingleResponse{}
		err = json.Unmarshal(data, &contents)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		if contents.Message != "success" {
			t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
		}

		if contents.Data.Name != testName {
			t.Errorf("Expected shelving unit name to be: %s. Got: %s.", testName, contents.Data.Name)
		}

		if contents.Data.Notes != testNotes {
			t.Errorf("Expected shelving unit notes to be: %s. Got: %s", testNotes, contents.Data.Notes)
		}
	})

	// Create Shelving Unit Test Case 2 - Valid input and does not include notes
	t.Run("CreateShelvingUnit-ValidNoNotes", func(t *testing.T) {
		testName := "Test Shelving Unit 2"

		values := map[string]string{"Name": testName}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateShelvingUnit(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		contents := shelvingUnitSingleResponse{}
		err = json.Unmarshal(data, &contents)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		if contents.Message != "success" {
			t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
		}

		if contents.Data.Name != testName {
			t.Errorf("Expected shelving unit name to be: %s. Got: %s.", testName, contents.Data.Name)
		}

		if contents.Data.Notes != "" {
			t.Errorf("Expected shelving unit notes to be: null. Got: %s", contents.Data.Notes)
		}
	})

	// Create Shelving Unit Test Case 3 - Invalid input, missing name
	t.Run("CreateShelvingUnit-InvalidMissingName", func(t *testing.T) {

		values := map[string]string{}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateShelvingUnit(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})
}
