package controllers

import (
	"bytes"
	"chi-boilerplate/infra/logger"
	"chi-boilerplate/repository"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type singleResponse struct {
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

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		logger.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}

func TestBuildingAPI(t *testing.T) {
	mockDB, _ := NewMockDB()
	handler := Handler{Repository: &repository.Repository{Database: mockDB}}
	r := chi.NewRouter()
	r.Post("/v1/entity-management/building", handler.CreateBuilding)

	srv := httptest.NewServer(r)
	defer srv.Close()

	t.Run("CreateBuilding-ValidAllFields", func(t *testing.T) {
		testName := "Test Building 1"
		testAddress := "Test Address 1"
		testNotes := "Test Notes 1"

		values := map[string]string{"Name": testName, "Address": testAddress, "Notes": testNotes}
		payload, _ := json.Marshal(values)

		req := httptest.NewRequest(http.MethodPost, srv.URL+"/v1/entity-management/building", bytes.NewBuffer(payload))
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

		contents := singleResponse{}
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
}
