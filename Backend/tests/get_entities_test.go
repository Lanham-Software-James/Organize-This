// Package tests is where all of out unit tests are described.
package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"organize-this/controllers"
	"organize-this/models"
	"organize-this/repository"
	"organize-this/routers"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type getEntitiesSingleResponse struct {
	Message string                     `json:"message"`
	Data    models.GetEntitiesResponse `json:"data"`
}

func setupMockExpectations(mock sqlmock.Sqlmock, getEntitiesQuery string) {
	mock.ExpectQuery(getEntitiesQuery).
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
}

func processResponse(t *testing.T, w *httptest.ResponseRecorder) (contents getEntitiesSingleResponse) {
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil. Got: %v", err)
	}

	contents = getEntitiesSingleResponse{}
	err = json.Unmarshal(data, &contents)

	return contents
}

func isValidResponse(t *testing.T, contents getEntitiesSingleResponse) {
	if contents.Message != "success" {
		t.Errorf("Expected message to be 'success'. Got: %s", contents.Message)
	}

	dataType := reflect.TypeOf(contents.Data).String()
	if dataType != "models.GetEntitiesResponse" {
		t.Errorf("Expected data to be type models.GetEntitiesResponse. Got: %v", dataType)
	}

	dataType = reflect.TypeOf(contents.Data.TotalCount).String()
	if dataType != "int" {
		t.Errorf("Expected TotalCount to be type int. Got: %v", dataType)
	}

	dataType = reflect.TypeOf(contents.Data.Entities).String()
	if dataType != "[]models.GetEntitiesResponseData" {
		t.Errorf("Expected TotalCount to be type []models.GetEntitiesResponseData. Got: %v", dataType)
	}
}

// TestCreateEntityAll runs our unit tests for the CreateEntity function that apply to all categories.
func TestGetEntities(t *testing.T) {
	mockDB, mock := NewMockDB()

	handler := controllers.Handler{Repository: &repository.Repository{Database: mockDB}}

	r := routers.SetupRoute()

	srv := httptest.NewServer(r)
	defer srv.Close()

	endpoint := "/v1/entities"
	getEntitiesQueryBase := `SELECT 'building' AS category, id, name, notes, ' ' as location FROM buildings
							UNION ALL
							SELECT 'room' AS category, id, name, notes, ' ' as location FROM rooms
							UNION ALL
							SELECT 'shelving_unit' AS category, id, name, notes, ' ' as location FROM shelving_units
							UNION ALL
							SELECT 'shelf' AS category, id, name, notes, ' ' as location FROM shelves
							UNION ALL
							SELECT 'container' AS category, id, name, notes, ' ' as location FROM containers
							UNION ALL
							SELECT 'item' AS category, id, name, notes, ' ' as location FROM items`

	t.Run("BEUT-18: Get Entities Valid URL Params", func(t *testing.T) {
		offset := "0"
		limit := "20"

		setupMockExpectations(mock, getEntitiesQueryBase+" OFFSET "+offset+" LIMIT "+limit)

		params := "?offset=" + offset + "&limit=" + limit
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)

		contents := processResponse(t, w)

		isValidResponse(t, contents)
	})

	t.Run("BEUT-19: Get Entities No Offset Param", func(t *testing.T) {
		limit := "20"

		setupMockExpectations(mock, getEntitiesQueryBase+" LIMIT "+limit)

		params := "?limit=20"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)

		contents := processResponse(t, w)

		isValidResponse(t, contents)
	})

	t.Run("BEUT-20: Get Entities No Limit Param", func(t *testing.T) {
		offset := "0"

		setupMockExpectations(mock, getEntitiesQueryBase+" OFFSET "+offset)

		params := "?offset=0"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusOK, res.StatusCode)
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected error to be nil. Got: %v", err)
		}

		contents := getEntitiesSingleResponse{}
		err = json.Unmarshal(data, &contents)

		isValidResponse(t, contents)
	})

	t.Run("BEUT-21: Get Entities No Params", func(t *testing.T) {
		setupMockExpectations(mock, getEntitiesQueryBase+" OFFSET 0 LIMIT 20")

		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)

		contents := processResponse(t, w)

		isValidResponse(t, contents)
	})

	t.Run("BEUT-22: Get Entities Non-Zero Offset", func(t *testing.T) {
		offset := "20"
		limit := "20"

		setupMockExpectations(mock, getEntitiesQueryBase+" OFFSET "+offset+" LIMIT "+limit)

		params := "?offset=20&limit=20"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)

		contents := processResponse(t, w)

		isValidResponse(t, contents)
	})

	t.Run("BEUT-23: Get Entities Invalid Offset Non-Integer", func(t *testing.T) {
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

	t.Run("BEUT-24: Get Entities Invalid Limit Non-Integer", func(t *testing.T) {
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

	t.Run("BEUT-25: Get Entities Invalid Params Non-Integer", func(t *testing.T) {
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

	t.Run("BEUT-26: Get Entities Invalid Offset Negative Integer", func(t *testing.T) {
		params := "?offset=-20&limit=20"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("BEUT-27: Get Entities Invalid Limit Negative Integer", func(t *testing.T) {
		params := "?offset=0&limit=-20"
		req := httptest.NewRequest(http.MethodGet, srv.URL+endpoint+params, nil)
		w := httptest.NewRecorder()

		handler.GetEntities(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code to be: %d. Got: %d.", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("BEUT-28: Get Entities Invalid Params Negative Integer", func(t *testing.T) {
		params := "?offset=-20&limit=-20"
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
