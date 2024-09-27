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

		mock.ExpectQuery(`SELECT 'building' AS category, id, name, notes, ' ' as location FROM buildings
							UNION ALL
							SELECT 'room' AS category, id, name, notes, ' ' as location FROM rooms
							UNION ALL
							SELECT 'shelving_unit' AS category, id, name, notes, ' ' as location FROM shelving_units
							UNION ALL
							SELECT 'shelf' AS category, id, name, notes, ' ' as location FROM shelves
							UNION ALL
							SELECT 'container' AS category, id, name, notes, ' ' as location FROM containers
							UNION ALL
							SELECT 'item' AS category, id, name, notes, ' ' as location FROM items
							OFFSET 0
							LIMIT 20`).
			WillReturnRows(sqlmock.NewRows([]string{"category", "id", "name", "notes", "location"}).
				AddRow("building", 1, "Building 1", "", " ").
				AddRow("building", 2, "Building 2", "", " ").
				AddRow("room", 1, "Room 1", "", " ").
				AddRow("room", 2, "Room 2", "", " ").
				AddRow("shelving_unit", 1, "Shelving Unit 1", "", " ").
				AddRow("shelving_unit", 2, "Shelving Unit 2", "", " ").
				AddRow("shelf", 1, "Shelf 1", "", " ").
				AddRow("shelf", 2, "Shelf 2", "", " ").
				AddRow("container", 1, "Container 1", "", " ").
				AddRow("container", 2, "Container 2", "", " ").
				AddRow("item", 1, "Item 1", "", " ").
				AddRow("item", 2, "Item 2", "", " "))

		mock.ExpectQuery(`SELECT
							(SELECT COUNT\(\*\) FROM buildings) \+
							(SELECT COUNT\(\*\) FROM rooms) \+
							(SELECT COUNT\(\*\) FROM shelving_units) \+
							(SELECT COUNT\(\*\) FROM shelves) \+
							(SELECT COUNT\(\*\) FROM containers) \+
							(SELECT COUNT\(\*\) FROM items)
						AS EntityCount`).
			WillReturnRows(sqlmock.NewRows([]string{"EntityCount"}).AddRow(6))

		params := "?offset=0&limit=20"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}
	})

	t.Run("BEUT-19: Get Entities No Offset Param", func(t *testing.T) {

		mock.ExpectQuery(`SELECT 'building' AS category, id, name, notes, ' ' as location FROM buildings
							UNION ALL
							SELECT 'room' AS category, id, name, notes, ' ' as location FROM rooms
							UNION ALL
							SELECT 'shelving_unit' AS category, id, name, notes, ' ' as location FROM shelving_units
							UNION ALL
							SELECT 'shelf' AS category, id, name, notes, ' ' as location FROM shelves
							UNION ALL
							SELECT 'container' AS category, id, name, notes, ' ' as location FROM containers
							UNION ALL
							SELECT 'item' AS category, id, name, notes, ' ' as location FROM items
							LIMIT 20`).
			WillReturnRows(sqlmock.NewRows([]string{"category", "id", "name", "notes", "location"}).
				AddRow("building", 1, "Building 1", "", " ").
				AddRow("building", 2, "Building 2", "", " ").
				AddRow("room", 1, "Room 1", "", " ").
				AddRow("room", 2, "Room 2", "", " ").
				AddRow("shelving_unit", 1, "Shelving Unit 1", "", " ").
				AddRow("shelving_unit", 2, "Shelving Unit 2", "", " ").
				AddRow("shelf", 1, "Shelf 1", "", " ").
				AddRow("shelf", 2, "Shelf 2", "", " ").
				AddRow("container", 1, "Container 1", "", " ").
				AddRow("container", 2, "Container 2", "", " ").
				AddRow("item", 1, "Item 1", "", " ").
				AddRow("item", 2, "Item 2", "", " "))

		mock.ExpectQuery(`SELECT
				(SELECT COUNT\(\*\) FROM buildings) \+
				(SELECT COUNT\(\*\) FROM rooms) \+
				(SELECT COUNT\(\*\) FROM shelving_units) \+
				(SELECT COUNT\(\*\) FROM shelves) \+
				(SELECT COUNT\(\*\) FROM containers) \+
				(SELECT COUNT\(\*\) FROM items)
			AS EntityCount`).
			WillReturnRows(sqlmock.NewRows([]string{"EntityCount"}).AddRow(6))

		params := "?limit=20"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}
	})

	t.Run("BEUT-20: Get Entities No Limit Param", func(t *testing.T) {

		mock.ExpectQuery(`SELECT 'building' AS category, id, name, notes, ' ' as location FROM buildings
							UNION ALL
							SELECT 'room' AS category, id, name, notes, ' ' as location FROM rooms
							UNION ALL
							SELECT 'shelving_unit' AS category, id, name, notes, ' ' as location FROM shelving_units
							UNION ALL
							SELECT 'shelf' AS category, id, name, notes, ' ' as location FROM shelves
							UNION ALL
							SELECT 'container' AS category, id, name, notes, ' ' as location FROM containers
							UNION ALL
							SELECT 'item' AS category, id, name, notes, ' ' as location FROM items
							LIMIT 20`).
			WillReturnRows(sqlmock.NewRows([]string{"category", "id", "name", "notes", "location"}).
				AddRow("building", 1, "Building 1", "", " ").
				AddRow("building", 2, "Building 2", "", " ").
				AddRow("room", 1, "Room 1", "", " ").
				AddRow("room", 2, "Room 2", "", " ").
				AddRow("shelving_unit", 1, "Shelving Unit 1", "", " ").
				AddRow("shelving_unit", 2, "Shelving Unit 2", "", " ").
				AddRow("shelf", 1, "Shelf 1", "", " ").
				AddRow("shelf", 2, "Shelf 2", "", " ").
				AddRow("container", 1, "Container 1", "", " ").
				AddRow("container", 2, "Container 2", "", " ").
				AddRow("item", 1, "Item 1", "", " ").
				AddRow("item", 2, "Item 2", "", " "))

		mock.ExpectQuery(`SELECT
							(SELECT COUNT\(\*\) FROM buildings) \+
							(SELECT COUNT\(\*\) FROM rooms) \+
							(SELECT COUNT\(\*\) FROM shelving_units) \+
							(SELECT COUNT\(\*\) FROM shelves) \+
							(SELECT COUNT\(\*\) FROM containers) \+
							(SELECT COUNT\(\*\) FROM items)
						AS EntityCount`).
			WillReturnRows(sqlmock.NewRows([]string{"EntityCount"}).AddRow(6))

		params := "?offset=0"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}
	})

	t.Run("BEUT-21: Get Entities No Params", func(t *testing.T) {

		mock.ExpectQuery(`SELECT 'building' AS category, id, name, notes, ' ' as location FROM buildings
							UNION ALL
							SELECT 'room' AS category, id, name, notes, ' ' as location FROM rooms
							UNION ALL
							SELECT 'shelving_unit' AS category, id, name, notes, ' ' as location FROM shelving_units
							UNION ALL
							SELECT 'shelf' AS category, id, name, notes, ' ' as location FROM shelves
							UNION ALL
							SELECT 'container' AS category, id, name, notes, ' ' as location FROM containers
							UNION ALL
							SELECT 'item' AS category, id, name, notes, ' ' as location FROM items
							OFFSET 0
							LIMIT 20`).
			WillReturnRows(sqlmock.NewRows([]string{"category", "id", "name", "notes", "location"}).
				AddRow("building", 1, "Building 1", "", " ").
				AddRow("building", 2, "Building 2", "", " ").
				AddRow("room", 1, "Room 1", "", " ").
				AddRow("room", 2, "Room 2", "", " ").
				AddRow("shelving_unit", 1, "Shelving Unit 1", "", " ").
				AddRow("shelving_unit", 2, "Shelving Unit 2", "", " ").
				AddRow("shelf", 1, "Shelf 1", "", " ").
				AddRow("shelf", 2, "Shelf 2", "", " ").
				AddRow("container", 1, "Container 1", "", " ").
				AddRow("container", 2, "Container 2", "", " ").
				AddRow("item", 1, "Item 1", "", " ").
				AddRow("item", 2, "Item 2", "", " "))

		mock.ExpectQuery(`SELECT
							(SELECT COUNT\(\*\) FROM buildings) \+
							(SELECT COUNT\(\*\) FROM rooms) \+
							(SELECT COUNT\(\*\) FROM shelving_units) \+
							(SELECT COUNT\(\*\) FROM shelves) \+
							(SELECT COUNT\(\*\) FROM containers) \+
							(SELECT COUNT\(\*\) FROM items)
						AS EntityCount`).
			WillReturnRows(sqlmock.NewRows([]string{"EntityCount"}).AddRow(6))

		params := ""
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}
	})

	t.Run("BEUT-22: Get Entities Non-Zero Offset", func(t *testing.T) {

		mock.ExpectQuery(`SELECT 'building' AS category, id, name, notes, ' ' as location FROM buildings
							UNION ALL
							SELECT 'room' AS category, id, name, notes, ' ' as location FROM rooms
							UNION ALL
							SELECT 'shelving_unit' AS category, id, name, notes, ' ' as location FROM shelving_units
							UNION ALL
							SELECT 'shelf' AS category, id, name, notes, ' ' as location FROM shelves
							UNION ALL
							SELECT 'container' AS category, id, name, notes, ' ' as location FROM containers
							UNION ALL
							SELECT 'item' AS category, id, name, notes, ' ' as location FROM items
							OFFSET 20
							LIMIT 20`).
			WillReturnRows(sqlmock.NewRows([]string{"category", "id", "name", "notes", "location"}).
				AddRow("building", 1, "Building 1", "", " ").
				AddRow("building", 2, "Building 2", "", " ").
				AddRow("room", 1, "Room 1", "", " ").
				AddRow("room", 2, "Room 2", "", " ").
				AddRow("shelving_unit", 1, "Shelving Unit 1", "", " ").
				AddRow("shelving_unit", 2, "Shelving Unit 2", "", " ").
				AddRow("shelf", 1, "Shelf 1", "", " ").
				AddRow("shelf", 2, "Shelf 2", "", " ").
				AddRow("container", 1, "Container 1", "", " ").
				AddRow("container", 2, "Container 2", "", " ").
				AddRow("item", 1, "Item 1", "", " ").
				AddRow("item", 2, "Item 2", "", " "))

		mock.ExpectQuery(`SELECT
							(SELECT COUNT\(\*\) FROM buildings) \+
							(SELECT COUNT\(\*\) FROM rooms) \+
							(SELECT COUNT\(\*\) FROM shelving_units) \+
							(SELECT COUNT\(\*\) FROM shelves) \+
							(SELECT COUNT\(\*\) FROM containers) \+
							(SELECT COUNT\(\*\) FROM items)
						AS EntityCount`).
			WillReturnRows(sqlmock.NewRows([]string{"EntityCount"}).AddRow(6))

		params := "?offset=20&limit=20"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
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
