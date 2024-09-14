package controllers

import (
	"chi-boilerplate/helpers"
	"chi-boilerplate/infra/logger"
	"chi-boilerplate/models"
	"chi-boilerplate/repository"
	"encoding/json"
	"net/http"
)

func CreateBuilding(w http.ResponseWriter, request *http.Request) {
	building := new(models.Building)
	err := json.NewDecoder(request.Body).Decode(&building)
	if err != nil {
		logger.Errorf("Error creating building: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	repository.Save(&building)
	helpers.SuccessResponse(w, &building)
}
