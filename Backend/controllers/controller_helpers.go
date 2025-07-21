package controllers

import (
	"net/http"
	"willowsuite-vault/helpers"
	"willowsuite-vault/infra/logger"
)

// logAndRespond logs an error message and sends a bad request response.
func logAndRespond(w http.ResponseWriter, message string, err error) {
	if err != nil {
		logger.Errorf("%s: %s", message, err)
	} else {
		logger.Errorf(message)
	}
	helpers.BadRequest(w, message)
}
