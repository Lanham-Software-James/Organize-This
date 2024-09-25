package controllers

import (
	"net/url"
	"organize-this/infra/logger"
	"organize-this/models"
	"strconv"
)

type getEntitiesResponse struct {
	Entities []getEntitiesResponseData
}

type getEntitiesResponseData struct {
	ID       uint
	Name     string
	Category string
	Location string
	Notes    *string
}

type getEntitiesIntermediateEntity struct {
	Category string
	Entity   models.Entity
}

func getEntitiesParseQueryParams(values url.Values) (int, int, error) {
	offsetString := values.Get("offset")
	limitString := values.Get("limit")

	if offsetString == "" {
		offsetString = "0"
	}

	if limitString == "" {
		limitString = "20"
	}

	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		logger.Errorf("Error converting offset to int: %v", err)
		return 0, 0, err
	}

	limit, err := strconv.Atoi(limitString)
	if err != nil {
		logger.Errorf("Error converting limit to int: %v", err)
		return 0, 0, err
	}

	return offset, limit, nil
}

func getEntitiesBuildResponse(entities []getEntitiesIntermediateEntity) (response []getEntitiesResponseData) {
	for _, entity := range entities {
		data := getEntitiesResponseData{
			ID:       uint(entity.Entity.ID),
			Name:     entity.Entity.Name,
			Category: entity.Category,
			Location: " ",
			Notes:    entity.Entity.Notes,
		}

		response = append(response, data)
	}

	return response
}

func (handler Handler) getEntitiesCountEntities() (sum int) {
	sum = 0

	buildings, err := handler.Repository.Count(&models.Building{})
	if err != nil {
		logger.Errorf("Error counting buildings: %v", err)
		return
	}

	sum += buildings
	rooms, err := handler.Repository.Count(&models.Room{})
	if err != nil {
		logger.Errorf("Error counting rooms: %v", err)
		return
	}

	sum += rooms
	units, err := handler.Repository.Count(&models.ShelvingUnit{})
	if err != nil {
		logger.Errorf("Error counting shelving units: %v", err)
		return
	}

	sum += units
	shelves, err := handler.Repository.Count(&models.Shelf{})
	if err != nil {
		logger.Errorf("Error counting shelves: %v", err)
		return
	}

	sum += shelves
	containers, err := handler.Repository.Count(&models.Container{})
	if err != nil {
		logger.Errorf("Error counting containers: %v", err)
		return
	}

	sum += containers
	items, err := handler.Repository.Count(&models.Shelf{})
	if err != nil {
		logger.Errorf("Error counting items: %v", err)
		return
	}

	sum += items

	return sum
}
