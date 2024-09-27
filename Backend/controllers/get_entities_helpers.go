package controllers

import (
	"net/url"
	"organize-this/infra/logger"
	"strconv"
)

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
