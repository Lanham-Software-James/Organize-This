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

type buildingSingleResponse struct {
	Message string `json:"message"`
	Data    struct {
		ID        uint   `json:"ID"`
		Name      string `json:"Name"`
		Address   string `json:"Address"`
		Notes     string `json:"Notes"`
		CreatedAt string `json:"CreatedAt"`
		UpdatedAt string `json:"UpdatedAt"`
	}
}

func TestCreateBuilding(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}
	endpoint := "/v1/entity-management/building"
	r := chi.NewRouter()
	r.Post(endpoint, handler.CreateBuilding)

	srv := httptest.NewServer(r)
	defer srv.Close()

	// Create Building Test Case 1 - Valid input and includes all fields
	t.Run("CreateBuilding-ValidAllFields", func(t *testing.T) {
		testName := "Test Building 1"
		testAddress := "Test Address 1"
		testNotes := "Test Notes 1"

		values := map[string]string{"Name": testName, "Address": testAddress, "Notes": testNotes}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateBuilding(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		contents := buildingSingleResponse{}
		err = json.Unmarshal(data, &contents)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		if contents.Message != "success" {
			t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
		}

		if contents.Data.Name != testName {
			t.Errorf("Expected building name to be: %s. Got: %s.", testName, contents.Data.Name)
		}

		if contents.Data.Address != testAddress {
			t.Errorf("Expected building address to be: %s. Got: %s", testAddress, contents.Data.Address)
		}

		if contents.Data.Notes != testNotes {
			t.Errorf("Expected building notes to be: %s. Got: %s", testNotes, contents.Data.Notes)
		}
	})

	// Create Building Test Case 2 - Valid input and does not include notes
	t.Run("CreateBuilding-ValidNoNotes", func(t *testing.T) {
		testName := "Test Building 2"
		testAddress := "Test Address 2"

		values := map[string]string{"Name": testName, "Address": testAddress}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateBuilding(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		contents := buildingSingleResponse{}
		err = json.Unmarshal(data, &contents)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		if contents.Message != "success" {
			t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
		}

		if contents.Data.Name != testName {
			t.Errorf("Expected building name to be: %s. Got: %s.", testName, contents.Data.Name)
		}

		if contents.Data.Address != testAddress {
			t.Errorf("Expected building address to be: %s. Got: %s", testAddress, contents.Data.Address)
		}

		if contents.Data.Notes != "" {
			t.Errorf("Expected building notes to be: null. Got: %s", contents.Data.Notes)
		}
	})

	// Create Building Test Case 3 - Valid input and does not include address
	t.Run("CreateBuilding-ValidNoAddress", func(t *testing.T) {
		testName := "Test Building 3"
		testNotes := "Test Notes 3"

		values := map[string]string{"Name": testName, "Notes": testNotes}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateBuilding(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		contents := buildingSingleResponse{}
		err = json.Unmarshal(data, &contents)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		if contents.Message != "success" {
			t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
		}

		if contents.Data.Name != testName {
			t.Errorf("Expected building name to be: %s. Got: %s.", testName, contents.Data.Name)
		}

		if contents.Data.Address != "" {
			t.Errorf("Expected building address to be: null. Got: %s", contents.Data.Address)
		}

		if contents.Data.Notes != testNotes {
			t.Errorf("Expected building notes to be: %s. Got: %s", testNotes, contents.Data.Notes)
		}
	})

	// Create Building Test Case 4 - Valid input and does not include address or notes
	t.Run("CreateBuilding-ValidNoAddressNoNotes", func(t *testing.T) {
		testName := "Test Building 4"

		values := map[string]string{"Name": testName}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateBuilding(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		contents := buildingSingleResponse{}
		err = json.Unmarshal(data, &contents)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		if contents.Message != "success" {
			t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
		}

		if contents.Data.Name != testName {
			t.Errorf("Expected building name to be: %s. Got: %s.", testName, contents.Data.Name)
		}

		if contents.Data.Address != "" {
			t.Errorf("Expected building address to be: null. Got: %s", contents.Data.Address)
		}

		if contents.Data.Notes != "" {
			t.Errorf("Expected building notes to be: null. Got: %s", contents.Data.Notes)
		}
	})

	// Create Building Test Case 5 - Invalid input, missing name
	t.Run("CreateBuilding-InvalidMissingName", func(t *testing.T) {

		values := map[string]string{}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+endpoint, bytes.NewBuffer(payload))
		w := httptest.NewRecorder()

		handler.CreateBuilding(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})
}
