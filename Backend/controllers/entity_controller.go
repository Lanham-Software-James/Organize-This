// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"organize-this/helpers"
	"organize-this/models"
)

// CreateEntity returns void but sends a success message or error message back to the client.
func (handler Handler) CreateEntity(w http.ResponseWriter, request *http.Request) {
	byteData, err := io.ReadAll(request.Body)
	if err != nil {
		logAndRespond(w, "Error parsing request", err)
		return
	}

	var parsedData map[string]string
	if err = json.Unmarshal(byteData, &parsedData); err != nil {
		logAndRespond(w, "Error parsing json", err)
		return
	}

	name, category := parsedData["name"], parsedData["category"]
	if name == "" {
		logAndRespond(w, "Missing name", nil)
		return
	}

	if category == "" {
		logAndRespond(w, "Missing category", nil)
		return
	}

	id, err := handler.createEntityByCategory(category, parsedData)
	if err != nil {
		logAndRespond(w, err.Error(), nil)
		return
	}

	helpers.SuccessResponse(w, &id)
}

// GetEntities return void, but sends a paginated list of all entities back to the client.
func (handler Handler) GetEntities(w http.ResponseWriter, request *http.Request) {
	var response getEntitiesResponse
	var entities []getEntitiesIntermediateEntity

	values := request.URL.Query()
	offset, limit, err := getEntitiesParseQueryParams(values)
	if err != nil {
		logAndRespond(w, "Error reading query parameters", err)
		return
	}

	// Get Buildings
	var buildings []models.Building
	handler.Repository.Get(&buildings, offset, limit)

	for _, building := range buildings {
		entities = append(entities, getEntitiesIntermediateEntity{Category: "Building", Entity: building.Entity})
	}

	// If length left get rooms
	if len(entities) < limit {

		buildingCount, err := handler.Repository.Count(&models.Building{})
		if err != nil {
			logAndRespond(w, "Error counting buildings", err)
			return
		}

		offset = offset - int(buildingCount)

		var rooms []models.Room

		handler.Repository.Get(&rooms, offset, limit-len(entities))

		for _, room := range rooms {
			entities = append(entities, getEntitiesIntermediateEntity{Category: "Room", Entity: room.Entity})
		}
	}

	// If length left get shelving units
	if len(entities) < limit {
		roomCount, err := handler.Repository.Count(&models.Room{})
		if err != nil {
			logAndRespond(w, "Error counting rooms", err)
			return
		}

		offset = offset - int(roomCount)

		var units []models.ShelvingUnit
		handler.Repository.Get(&units, offset, limit-len(entities))

		for _, unit := range units {
			entities = append(entities, getEntitiesIntermediateEntity{Category: "Shelving Unit", Entity: unit.Entity})
		}
	}

	// If length left get shelves
	if len(entities) < limit {
		unitCount, err := handler.Repository.Count(&models.ShelvingUnit{})
		if err != nil {
			logAndRespond(w, "Error counting shelving units", err)
			return
		}

		offset = offset - int(unitCount)

		var shelves []models.Shelf
		handler.Repository.Get(&shelves, offset, limit-len(entities))

		for _, shelf := range shelves {
			entities = append(entities, getEntitiesIntermediateEntity{Category: "Shelf", Entity: shelf.Entity})
		}
	}

	// If length left get containers
	if len(entities) < limit {
		shelfCount, err := handler.Repository.Count(&models.Shelf{})
		if err != nil {
			logAndRespond(w, "Error counting shelves", err)
			return
		}

		offset = offset - int(shelfCount)

		var containers []models.Container
		handler.Repository.Get(&containers, offset, limit-len(entities))

		for _, container := range containers {
			entities = append(entities, getEntitiesIntermediateEntity{Category: "Container", Entity: container.Entity})
		}
	}

	// If length left get items
	if len(entities) < limit {
		containerCount, err := handler.Repository.Count(&models.Container{})
		if err != nil {
			logAndRespond(w, "Error counting containers", err)
			return
		}

		offset = offset - int(containerCount)

		var items []models.Item
		handler.Repository.Get(&items, offset, limit-len(entities))

		for _, item := range items {
			entities = append(entities, getEntitiesIntermediateEntity{Category: "Item", Entity: item.Entity})
		}
	}

	response = getEntitiesBuildResponse(entities)

	helpers.SuccessResponse(w, &response)
}
