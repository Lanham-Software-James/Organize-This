// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"net/http"
	"organize-this/helpers"
	"organize-this/infra/logger"
	"organize-this/models"
)

// CreateContainer returns void, but sends an http response with the newly created Container back to the client.
func (handler Handler) CreateContainer(w http.ResponseWriter, request *http.Request) {

	container := new(models.Container)

	err := json.NewDecoder(request.Body).Decode(&container)

	if err != nil {
		logger.Errorf("Error creating container: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if container.Name == "" {
		logger.Errorf("Error creating container: Missing name.")
		http.Error(w, "Missing container name.", http.StatusBadRequest)
		return
	}

	handler.Repository.Save(&container)
	helpers.SuccessResponse(w, &container)
}
