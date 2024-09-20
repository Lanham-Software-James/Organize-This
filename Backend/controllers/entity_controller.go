// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"organize-this/helpers"
	"organize-this/infra/logger"
)

// CreateEntity returns void, but sends an success message or error message back to the client
func (handler Handler) CreateEntity(w http.ResponseWriter, request *http.Request) {

	var id uint

	byteData, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Errorf("Error parsing request: %s", err)
		helpers.BadRequest(w, err)
		return
	}

	parsedData := map[string]string{}

	err = json.Unmarshal(byteData, &parsedData)
	if err != nil {
		logger.Errorf("Error parsing json: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if parsedData["name"] == "" {
		logger.Errorf("Error creating entity: Missing name.")
		helpers.BadRequest(w, "Missing name")
		return
	}

	switch parsedData["category"] {
	case "item":
		id = handler.addItem(parsedData)
		break

	case "container":
		id = handler.addContainer(parsedData)
		break

	case "shelf":
		id = handler.addShelf(parsedData)
		break

	case "shelvingunit":
		id = handler.addShelvingUnit(parsedData)
		break

	case "room":
		id = handler.addRoom(parsedData)
		break

	case "building":
		id = handler.addBuilding(parsedData)
		break

	default:
		logger.Errorf("Error creating entity: Invalid category")
		helpers.BadRequest(w, "Invalid category")
		return
	}

	helpers.SuccessResponse(w, &id)
}
