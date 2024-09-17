// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"net/http"
	"organize-this/helpers"
	"organize-this/infra/logger"
	"organize-this/models"
)

// CreateItem returns void, but sends an http response with the newly created Item back to the client.
func (handler Handler) CreateItem(w http.ResponseWriter, request *http.Request) {

	item := new(models.Item)

	err := json.NewDecoder(request.Body).Decode(&item)

	if err != nil {
		logger.Errorf("Error creating item: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if item.Name == "" {
		logger.Errorf("Error creating item: Missing name.")
		http.Error(w, "Missing item name.", http.StatusBadRequest)
		return
	}

	handler.Repository.Save(&item)
	helpers.SuccessResponse(w, &item)
}
