// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"net/http"
	"organize-this/helpers"
	"organize-this/infra/logger"
	"organize-this/models"
)

// CreateBuilding returns void, but sends an http response with the newly created Building back to the client.
func (handler Handler) CreateBuilding(w http.ResponseWriter, request *http.Request) {

	building := new(models.Building)

	err := json.NewDecoder(request.Body).Decode(&building)

	if err != nil {
		logger.Errorf("Error creating building: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if building.Name == "" {
		logger.Errorf("Error creating building: Missing name.")
		http.Error(w, "Missing building name.", http.StatusBadRequest)
		return
	}

	handler.Repository.Save(&building)
	helpers.SuccessResponse(w, &building)
}
