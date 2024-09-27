// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"organize-this/helpers"
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

	id, err := handler.createEntityByCategory(category, parsedData)
	if err != nil {
		logAndRespond(w, err.Error(), nil)
		return
	}

	helpers.SuccessResponse(w, &id)
}

// GetEntities return void, but sends a paginated list of all entities back to the client.
func (handler Handler) GetEntities(w http.ResponseWriter, request *http.Request) {
	var response getEntitiesResponse

	values := request.URL.Query()
	offset, limit, err := getEntitiesParseQueryParams(values)
	if err != nil {
		logAndRespond(w, "Error reading query parameters", err)
		return
	}

	response.Entities = handler.Repository.GetAllEntities(offset, limit)

	response.TotalCount = handler.Repository.CountEntities()
	helpers.SuccessResponse(w, &response)
}
