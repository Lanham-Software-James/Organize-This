// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"organize-this/helpers"
	"organize-this/infra/logger"
	"organize-this/models"
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

// GetEntites return void, but sends a paginated list of all entities back to the client.
func (handler Handler) GetEntities(w http.ResponseWriter, request *http.Request) {
	var response helpers.GetEntities_Response
	var entities []helpers.GetEntities_IntermediateEntity

	values := request.URL.Query()
	offset, limit, err := helpers.GetEntities_ParseQueryParams(values)
	if err != nil {
		helpers.BadRequest(w, err)
		return
	}

	// Get Buildings
	var buildings []models.Building
	handler.Repository.Get(&buildings, offset, limit)

	for _, building := range buildings {
		entities = append(entities, helpers.GetEntities_IntermediateEntity{Category: "Building", Entity: building.Entity})
	}

	// If length left get rooms
	if len(entities) < limit {

		buildingCount, _ := handler.Repository.Count(&models.Building{})
		offset = offset - int(buildingCount)

		var rooms []models.Room

		handler.Repository.Get(&rooms, offset, limit-len(entities))

		for _, room := range rooms {
			entities = append(entities, helpers.GetEntities_IntermediateEntity{Category: "Room", Entity: room.Entity})
		}
	}

	// If length left get shelving units
	if len(entities) < limit {
		roomCount, _ := handler.Repository.Count(&models.Room{})

		offset = offset - int(roomCount)

		var units []models.ShelvingUnit
		handler.Repository.Get(&units, offset, limit-len(entities))

		for _, unit := range units {
			entities = append(entities, helpers.GetEntities_IntermediateEntity{Category: "Shelving Unit", Entity: unit.Entity})
		}
	}

	// If length left get shelves
	if len(entities) < limit {
		unitCount, _ := handler.Repository.Count(&models.ShelvingUnit{})

		offset = offset - int(unitCount)

		var shelves []models.Shelf
		handler.Repository.Get(&shelves, offset, limit-len(entities))

		for _, shelf := range shelves {
			entities = append(entities, helpers.GetEntities_IntermediateEntity{Category: "Shelf", Entity: shelf.Entity})
		}
	}

	// If length left get containers
	if len(entities) < limit {
		shelfCount, _ := handler.Repository.Count(&models.Shelf{})

		offset = offset - int(shelfCount)

		var containers []models.Shelf
		handler.Repository.Get(&containers, offset, limit-len(entities))

		for _, container := range containers {
			entities = append(entities, helpers.GetEntities_IntermediateEntity{Category: "Container", Entity: container.Entity})
		}
	}

	// If length left get items
	if len(entities) < limit {
		containerCount, _ := handler.Repository.Count(&models.Container{})

		offset = offset - int(containerCount)

		var items []models.Item
		handler.Repository.Get(&items, offset, limit-len(entities))

		for _, item := range items {
			entities = append(entities, helpers.GetEntities_IntermediateEntity{Category: "Item", Entity: item.Entity})
		}
	}

	response = helpers.GetEntities_BuildResponse(entities)

	helpers.SuccessResponse(w, &response)
}
