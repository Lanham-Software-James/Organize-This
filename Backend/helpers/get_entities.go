package helpers

import (
	"net/url"
	"organize-this/infra/logger"
	"organize-this/models"
	"strconv"
)

type GetEntities_Response struct {
	Count    int
	Entities []GetEntities_ResponseData
}

type GetEntities_ResponseData struct {
	ID       uint
	Name     string
	Category string
	Location string
	Notes    *string
}

type GetEntities_IntermediateEntity struct {
	Category string
	Entity   models.Entity
}

func GetEntities_ParseQueryParams(values url.Values) (int, int, error) {
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

func GetEntities_BuildResponse(entities []GetEntities_IntermediateEntity) (response GetEntities_Response) {
	response = GetEntities_Response{
		Count:    len(entities),
		Entities: []GetEntities_ResponseData{},
	}

	for _, entity := range entities {
		data := GetEntities_ResponseData{
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
