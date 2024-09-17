package controllers

import (
	"chi-boilerplate/helpers"
	"chi-boilerplate/infra/logger"
	"chi-boilerplate/models"
	"encoding/json"
	"net/http"
)

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
		http.Error(w, "Missing room name.", http.StatusBadRequest)
		return
	}

	handler.Repository.Save(&room)
	helpers.SuccessResponse(w, &room)
}
