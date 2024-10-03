// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"organize-this/helpers"
	"organize-this/models"

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

	id, err := handler.createEntityByCategory(claims["username"].(string), category, parsedData)
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
	response.Entities = handler.Repository.GetAllEntities(userID, offset, limit)

	response.TotalCount = handler.Repository.CountEntities(userID)
	helpers.SuccessResponse(w, &response)
}
