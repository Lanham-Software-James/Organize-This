// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"organize-this/helpers"
	"organize-this/models"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

// CreateEntity returns void but sends a success message or error message back to the client.
func (handler Handler) CreateEntity(w http.ResponseWriter, request *http.Request) {
	byteData, err := io.ReadAll(request.Body)
	if err != nil {
		logAndRespond(w, "Error parsing request", err)
		return
	}

	var parsedData map[string]string
	if err = json.Unmarshal(byteData, &parsedData); err != nil {
		logAndRespond(w, "Error parsing json", err)
		return
	}

	name, category := parsedData["name"], parsedData["category"]
	if name == "" {
		logAndRespond(w, "Missing name", nil)
		return
	}

	if category == "" {
		logAndRespond(w, "Missing category", nil)
		return
	}

	claims := request.Context().Value("user_claims").(jwt.MapClaims)

	id, err := handler.createEntityByCategory(request.Context(), claims["username"].(string), category, parsedData)
	if err != nil {
		logAndRespond(w, err.Error(), nil)
		return
	}

	helpers.SuccessResponse(w, &id)
}

// GetEntities return void, but sends a paginated list of all entities back to the client.
func (handler Handler) GetEntities(w http.ResponseWriter, request *http.Request) {
	var response models.GetEntitiesResponse

	values := request.URL.Query()
	offset, limit, err := getEntitiesParseQueryParams(values)
	if err != nil {
		logAndRespond(w, "Error reading query parameters", err)
		return
	}

	claims := request.Context().Value("user_claims").(jwt.MapClaims)
	userID := claims["username"].(string)
	response.Entities = handler.Repository.GetAllEntities(request.Context(), userID, offset, limit)

	response.TotalCount = handler.Repository.CountEntities(request.Context(), userID)
	helpers.SuccessResponse(w, &response)
}

// GetEntity return void, but sends a single entity back to the client if it finds a match.
func (handler Handler) GetEntity(w http.ResponseWriter, request *http.Request) {
	category := chi.URLParam(request, "category")
	idParam := chi.URLParam(request, "id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logAndRespond(w, fmt.Sprintf("ID must be type integer: %v", idParam), nil)
	}

	claims := request.Context().Value("user_claims").(jwt.MapClaims)
	userID := claims["username"].(string)

	var model interface{}

	switch category {
	case "item":
		model = &models.Item{
			Entity: models.Entity{
				ID: id,
			},
		}
	case "container":
		model = &models.Container{
			Entity: models.Entity{
				ID: id,
			},
		}
	case "shelf":
		model = &models.Shelf{
			Entity: models.Entity{
				ID: id,
			},
		}
	case "shelvingunit":
		model = &models.ShelvingUnit{
			Entity: models.Entity{
				ID: id,
			},
		}
	case "room":
		model = &models.Room{
			Entity: models.Entity{
				ID: id,
			},
		}
	case "building":
		model = &models.Building{
			Entity: models.Entity{
				ID: id,
			},
		}
	default:
		logAndRespond(w, fmt.Sprintf("Invalid Category: %v", category), nil)
	}

	dberr := handler.Repository.GetOne(model, userID)
	if dberr != nil {
		logAndRespond(w, fmt.Sprintf("Entity category of %v with id %v not found.", category, id), nil)
		return
	}

	helpers.SuccessResponse(w, model)
}
