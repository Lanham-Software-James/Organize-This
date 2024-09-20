// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"net/http"
	"organize-this/helpers"
	"organize-this/infra/logger"
	"organize-this/models"
)

// CreateRoom returns void, but sends an http response with the newly created Room back to the client.
func (handler Handler) CreateRoom(w http.ResponseWriter, request *http.Request) {

	room := new(models.Room)

	err := json.NewDecoder(request.Body).Decode(&room)

	if err != nil {
		logger.Errorf("Error creating room: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if room.Name == "" {
		logger.Errorf("Error creating room: Missing name.")
		helpers.BadRequest(w, "Missing Building Name.")
		return
	}

	handler.Repository.Save(&room)
	helpers.SuccessResponse(w, &room)
}
