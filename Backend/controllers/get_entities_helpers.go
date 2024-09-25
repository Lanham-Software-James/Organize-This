package controllers

import (
	"net/url"
	"organize-this/infra/logger"
	"organize-this/models"
	"strconv"
)

type getEntitiesResponse struct {
	Count    int
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

func getEntitiesBuildResponse(entities []getEntitiesIntermediateEntity) (response getEntitiesResponse) {
	response = getEntitiesResponse{
		Count:    len(entities),
		Entities: []getEntitiesResponseData{},
	}

	for _, entity := range entities {
		data := getEntitiesResponseData{
			ID:       uint(entity.Entity.ID),
			Name:     entity.Entity.Name,
			Category: entity.Category,
			Location: " ",
			Notes:    entity.Entity.Notes,
		}

		response.Entities = append(response.Entities, data)
	}

	return response
}
