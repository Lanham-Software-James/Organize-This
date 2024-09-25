// Package tests is where all of out unit tests are described.
package tests

import (
	"net/http"
	"net/http/httptest"
	"organize-this/controllers"
	"organize-this/repository"
	"organize-this/routers"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestCreateEntityAll runs our unit tests for the CreateEntity function that apply to all categories.
func TestGetEntities(t *testing.T) {
	mockDB, mock := NewMockDB()

	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}

	r := routers.SetupRoute()

	srv := httptest.NewServer(r)
	defer srv.Close()

	endpoint := "/v1/entities"

	t.Run("BEUT-18: Get Entities Valid URL Params", func(t *testing.T) {

		mock.ExpectQuery(`SELECT \* FROM "buildings" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Building 1").AddRow(2, "Building 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "buildings"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "rooms" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Room 1").AddRow(2, "Room 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "rooms"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "shelving_units" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Shelving Unit 1").AddRow(2, "Shelving Unit 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "shelving_units"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "shelves" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Shelf 1").AddRow(2, "Shelf 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "shelves"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "containers" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Container 1").AddRow(2, "Container 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "containers"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "items" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Item 1").AddRow(2, "Item 2"))

		params := "?offset=0&limit=20"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		contentRange := res.Header.Get("Content-Range")
		if contentRange != "0-20/0" {
			t.Errorf("Expected Content-Range header to be: 0-20/0. Got: %v.", contentRange)
		}
	})

	t.Run("BEUT-19: Get Entities No Offset Param", func(t *testing.T) {

		mock.ExpectQuery(`SELECT \* FROM "buildings" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Building 1").AddRow(2, "Building 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "buildings"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "rooms" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Room 1").AddRow(2, "Room 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "rooms"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "shelving_units" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Shelving Unit 1").AddRow(2, "Shelving Unit 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "shelving_units"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "shelves" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Shelf 1").AddRow(2, "Shelf 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "shelves"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "containers" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Container 1").AddRow(2, "Container 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "containers"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "items" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Item 1").AddRow(2, "Item 2"))

		params := "?limit=20"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		contentRange := res.Header.Get("Content-Range")
		if contentRange != "0-20/0" {
			t.Errorf("Expected Content-Range header to be: 0-20/0. Got: %v.", contentRange)
		}
	})

	t.Run("BEUT-20: Get Entities No Limit Param", func(t *testing.T) {

		mock.ExpectQuery(`SELECT \* FROM "buildings" LIMIT 20`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Building 1").AddRow(2, "Building 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "buildings"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "rooms" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Room 1").AddRow(2, "Room 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "rooms"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "shelving_units" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Shelving Unit 1").AddRow(2, "Shelving Unit 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "shelving_units"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "shelves" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Shelf 1").AddRow(2, "Shelf 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "shelves"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "containers" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Container 1").AddRow(2, "Container 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "containers"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "items" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Item 1").AddRow(2, "Item 2"))

		params := "?offset=0"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		contentRange := res.Header.Get("Content-Range")
		if contentRange != "0-20/0" {
			t.Errorf("Expected Content-Range header to be: 0-20/0. Got: %v.", contentRange)
		}
	})

	t.Run("BEUT-21: Get Entities No Params", func(t *testing.T) {

		mock.ExpectQuery(`SELECT \* FROM "buildings" LIMIT 20`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Building 1").AddRow(2, "Building 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "buildings"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "rooms" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Room 1").AddRow(2, "Room 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "rooms"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "shelving_units" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Shelving Unit 1").AddRow(2, "Shelving Unit 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "shelving_units"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "shelves" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Shelf 1").AddRow(2, "Shelf 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "shelves"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "containers" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Container 1").AddRow(2, "Container 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "containers"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "items" LIMIT \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Item 1").AddRow(2, "Item 2"))

		params := ""
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		contentRange := res.Header.Get("Content-Range")
		if contentRange != "0-20/0" {
			t.Errorf("Expected Content-Range header to be: 0-20/0. Got: %v.", contentRange)
		}
	})

	t.Run("BEUT-22: Get Entities Non-Zero Offset", func(t *testing.T) {

		mock.ExpectQuery(`SELECT \* FROM "buildings" LIMIT \d+ OFFSET 20`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Building 1").AddRow(2, "Building 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "buildings"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "rooms" LIMIT \d+ OFFSET \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Room 1").AddRow(2, "Room 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "rooms"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "shelving_units" LIMIT \d+ OFFSET \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Shelving Unit 1").AddRow(2, "Shelving Unit 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "shelving_units"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "shelves" LIMIT \d+ OFFSET \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Shelf 1").AddRow(2, "Shelf 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "shelves"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "containers" LIMIT \d+ OFFSET \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Container 1").AddRow(2, "Container 2"))
		mock.ExpectQuery(`SELECT count\(\*\) FROM "containers"`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery(`SELECT \* FROM "items" LIMIT \d+ OFFSET \d+`).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Item 1").AddRow(2, "Item 2"))

		params := "?offset=20&limit=20"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		contentRange := res.Header.Get("Content-Range")
		if contentRange != "20-40/0" {
			t.Errorf("Expected Content-Range header to be: 20-40/0. Got: %v.", contentRange)
		}
	})

	t.Run("BEUT-23: Get Entities Invalid Offset", func(t *testing.T) {

		params := "?offset=test&limit=20"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("BEUT-24: Get Entities Invalid Limit", func(t *testing.T) {

		params := "?offset=0&limit=test"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("BEUT-24: Get Entities Invalid Params", func(t *testing.T) {

		params := "?offset=test&limit=tests"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})
}
