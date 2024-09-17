// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"net/http"
	"organize-this/helpers"
	"organize-this/infra/logger"
	"organize-this/models"
)

// CreateShelf returns void, but sends an http response with the newly created Shelf back to the client.
func (handler Handler) CreateShelf(w http.ResponseWriter, request *http.Request) {

	shelf := new(models.Shelf)

	err := json.NewDecoder(request.Body).Decode(&shelf)

	if err != nil {
		logger.Errorf("Error creating shelf: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if shelf.Name == "" {
		logger.Errorf("Error creating shelf: Missing name.")
		http.Error(w, "Missing shelf name.", http.StatusBadRequest)
		return
	}

	handler.Repository.Save(&shelf)
	helpers.SuccessResponse(w, &shelf)
}
