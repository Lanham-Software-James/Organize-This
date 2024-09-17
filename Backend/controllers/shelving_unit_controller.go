package controllers

import (
	"encoding/json"
	"net/http"
	"organize-this/helpers"
	"organize-this/infra/logger"
	"organize-this/models"
)

func (handler Handler) CreateShelvingUnit(w http.ResponseWriter, request *http.Request) {

	shelvingUnit := new(models.ShelvingUnit)

	err := json.NewDecoder(request.Body).Decode(&shelvingUnit)

	if err != nil {
		logger.Errorf("Error creating shelving unit: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if shelvingUnit.Name == "" {
		logger.Errorf("Error creating shelving unit: Missing name.")
		http.Error(w, "Missing shelving unit name.", http.StatusBadRequest)
		return
	}

	handler.Repository.Save(&shelvingUnit)
	helpers.SuccessResponse(w, &shelvingUnit)
}
