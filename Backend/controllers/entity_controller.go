// Package controllers provides all the various controllers for requests coming to the Organize This API.
package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"organize-this/helpers"
	"organize-this/infra/logger"
	"organize-this/models"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

// CreateEntity returns void but sends a success message or error message back to the client.
func (handler Handler) CreateEntity(w http.ResponseWriter, request *http.Request) {
	byteData, err := io.ReadAll(request.Body)
	if err != nil {
		logAndRespond(w, "Error parsing request", err)
		return
	}

	var parsedData map[string]string
	err = json.Unmarshal(byteData, &parsedData)
	if err != nil {
		logAndRespond(w, "Error parsing json", err)
		return
	}

	_, name, category, parentID, parentCategory, err := validateParams(parsedData, false)
	if err != nil {
		logAndRespond(w, "Error validating parameters", err)
		return
	}

	claims := request.Context().Value("user_claims").(jwt.MapClaims)
	userID := claims["username"].(string)
	tmpNotes := parsedData["notes"]

	var parent models.Parent
	if category != "building" {
		validParent := false
		validParent, parent = buildParent(category, parentID, parentCategory)
		if !validParent {
			logAndRespond(w, "Invalid parent.", nil)
			return
		}
	}

	entity := models.Entity{
		Name:   name,
		Notes:  &tmpNotes,
		UserID: userID,
	}

	validEntity, model := buildEntity(entity, parent, category, parsedData["address"])
	if !validEntity {
		logAndRespond(w, fmt.Sprintf("Invalid category %v.", category), nil)
		return
	}

	dberr := handler.Repository.Save(model)
	if dberr != nil {
		logAndRespond(w, "Error adding etity.", nil)
		return
	}

	handler.Repository.FlushEntities(request.Context(), userID)
	helpers.SuccessResponse(w, &model)
}

// GetEntities return void, but sends a paginated list of all entities back to the client.
func (handler Handler) GetEntities(w http.ResponseWriter, request *http.Request) {
	var response models.GetEntitiesResponse

	values := request.URL.Query()
	offset, limit, search, filters, err := getEntitiesParseQueryParams(values)
	if err != nil {
		logAndRespond(w, "Error reading query parameters", err)
		return
	}

	claims := request.Context().Value("user_claims").(jwt.MapClaims)
	userID := claims["username"].(string)
	response.Entities, _ = handler.Repository.GetAllEntities(request.Context(), userID, offset, limit, search, filters)

	response.TotalCount = handler.Repository.CountEntities(request.Context(), userID, search, filters)
	helpers.SuccessResponse(w, &response)
}

// GetEntity return void, but sends a single entity back to the client if it finds a match.
func (handler Handler) GetEntity(w http.ResponseWriter, request *http.Request) {
	category := chi.URLParam(request, "category")
	idParam := chi.URLParam(request, "id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logAndRespond(w, fmt.Sprintf("ID must be type integer: %v", idParam), nil)
		return
	}

	claims := request.Context().Value("user_claims").(jwt.MapClaims)
	userID := claims["username"].(string)

	entity := models.Entity{
		ID: id,
	}

	validEntity, model := buildEntity(entity, models.Parent{}, category, "")
	if !validEntity {
		logAndRespond(w, fmt.Sprintf("Invalid category %v.", category), nil)
		return
	}

	dberr := handler.Repository.GetOne(model, userID)
	if dberr != nil {
		logAndRespond(w, fmt.Sprintf("Entity category of %v with id %v not found.", category, id), nil)
		return
	}

	helpers.SuccessResponse(w, model)
}

// EditEntity returns void, but sends a success message or error message back to the client.
func (handler Handler) EditEntity(w http.ResponseWriter, request *http.Request) {
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

	id, name, category, parentID, parentCategory, err := validateParams(parsedData, true)
	if err != nil {
		logAndRespond(w, "Error validating parameters: ", err)
		return
	}

	claims := request.Context().Value("user_claims").(jwt.MapClaims)
	userID := claims["username"].(string)
	tmpNotes := parsedData["notes"]

	var parent models.Parent
	if category != "building" {
		validParent := false
		validParent, parent = buildParent(category, parentID, parentCategory)
		if !validParent {
			logAndRespond(w, "Invalid parent.", nil)
			return
		}
	}

	entity := models.Entity{
		ID:     id,
		Name:   name,
		Notes:  &tmpNotes,
		UserID: userID,
	}

	validEntity, model := buildEntity(entity, parent, category, parsedData["address"])
	if !validEntity {
		logAndRespond(w, fmt.Sprintf("Invalid category %v.", category), nil)
		return
	}

	dberr := handler.Repository.Save(model)
	if dberr != nil {
		logAndRespond(w, fmt.Sprintf("Entity category of %v with id %v not found.", category, id), nil)
		return
	}

	handler.Repository.FlushEntities(request.Context(), userID)
	helpers.SuccessResponse(w, model)
}

// DeleteEntity return void, but sends a confirmation message to the client.
func (handler Handler) DeleteEntity(w http.ResponseWriter, request *http.Request) {
	category := chi.URLParam(request, "category")
	idParam := chi.URLParam(request, "id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		logAndRespond(w, fmt.Sprintf("ID must be type integer: %v", idParam), nil)
		return
	}

	claims := request.Context().Value("user_claims").(jwt.MapClaims)
	userID := claims["username"].(string)

	entity := models.Entity{
		ID: id,
	}

	validEntity, model := buildEntity(entity, models.Parent{}, category, "")
	if !validEntity {
		logAndRespond(w, fmt.Sprintf("Invalid category %v.", category), nil)
		return
	}

	dberr := handler.Repository.GetOne(model, userID)
	if dberr != nil {
		logAndRespond(w, fmt.Sprintf("Entity category of %v with id %v not found.", category, id), nil)
		return
	}

	if category != "item" {
		hasChildren, count, err := handler.Repository.HasChildren(entity.ID, category, userID)
		if err != nil {
			logAndRespond(w, "Issue getting children", err)
			return
		}

		if hasChildren {
			logAndRespond(w,
				fmt.Sprintf("Cannot delete entity with children. Number of children: %d", count),
				fmt.Errorf("Cannot delete entity with children. Number of children: %v", count))
			return
		}
	}

	dberr = handler.Repository.Delete(model, userID)
	if dberr != nil {
		logAndRespond(w,
			fmt.Sprintf("Error deleting entity: %s - %d", category, id),
			fmt.Errorf("Error deleting entity: %s - %d", category, id))
		return
	}

	handler.Repository.FlushEntities(request.Context(), userID)
	helpers.SuccessResponse(w, "Successfully Deleted!")
}

// GetParents returns void, but sends valid parents back to the client.
func (handler Handler) GetParents(w http.ResponseWriter, request *http.Request) {
	category := chi.URLParam(request, "category")

	if category != "item" &&
		category != "container" &&
		category != "shelf" &&
		category != "shelving_unit" &&
		category != "room" {
		logAndRespond(w, "Invalid category", nil)
		return
	}

	claims := request.Context().Value("user_claims").(jwt.MapClaims)
	userID := claims["username"].(string)

	parents, err := handler.Repository.GetParents(request.Context(), category, userID)
	if err != nil {
		logAndRespond(w, "Issue getting parents.", err)
		return
	}

	helpers.SuccessResponse(w, parents)
}

func validateParams(parsedData map[string]string, edit bool) (uint64, string, string, uint64, string, error) {
	var err error
	var id uint64
	if parsedData["id"] == "" && edit {
		err = errors.New("Missing id")
	}

	if parsedData["name"] == "" {
		err = errors.New("Missing name")
	}

	if parsedData["category"] == "" {
		err = errors.New("Missing category")
	}

	if parsedData["parentID"] == "" && parsedData["category"] != "building" {
		err = errors.New("Missing parent id")
	}

	if parsedData["parentCategory"] == "" && parsedData["category"] != "building" {
		err = errors.New("Missing parent category")
	}

	if edit {
		id, err = strconv.ParseUint(parsedData["id"], 10, 64)
		if err != nil {
			err = errors.New("ID must be type integer")
		}
	}

	parentID, err2 := strconv.ParseUint(parsedData["parentID"], 10, 64)
	if err2 != nil && parsedData["category"] != "building" {
		err = errors.New("Parent ID must be type integer")
	}

	return id, parsedData["name"], parsedData["category"], parentID, parsedData["parentCategory"], err
}

func validateParent(category string, parentCategory string) bool {
	isValid := false

	switch category {
	case "item":
		if (parentCategory == "container") || (parentCategory == "shelf") || (parentCategory == "room") {
			isValid = true
		}
		break
	case "container":
		if (parentCategory == "shelf") || (parentCategory == "room") {
			isValid = true
		}
		break
	case "shelf":
		if parentCategory == "shelving_unit" {
			isValid = true
		}
		break
	case "shelving_unit":
		if parentCategory == "room" {
			isValid = true
		}
		break
	case "room":
		if parentCategory == "building" {
			isValid = true
		}
		break
	default:
		logger.Errorf("Invalid category for entity.")
	}

	return isValid
}

func buildParent(category string, parentID uint64, parentCategory string) (bool, models.Parent) {
	var parent models.Parent

	isParentValid := validateParent(category, parentCategory)

	if isParentValid {
		parent = models.Parent{
			ParentID:       parentID,
			ParentCategory: parentCategory,
		}
	}

	return isParentValid, parent
}

func buildEntity(entity models.Entity, parent models.Parent, category string, address string) (bool, interface{}) {
	valid := true
	var model interface{}

	switch category {
	case "item":
		model = &models.Item{
			Entity: entity,
			Parent: parent,
		}
		break
	case "container":
		model = &models.Container{
			Entity: entity,
			Parent: parent,
		}
		break
	case "shelf":
		model = &models.Shelf{
			Entity: entity,
			Parent: parent,
		}
		break
	case "shelving_unit":
		model = &models.ShelvingUnit{
			Entity: entity,
			Parent: parent,
		}
		break
	case "room":
		model = &models.Room{
			Entity: entity,
			Parent: parent,
		}
		break
	case "building":
		tmpAddress := address

		model = &models.Building{
			Entity:  entity,
			Address: &tmpAddress,
		}
		break
	default:
		logger.Errorf("Invalid Category: %v", category)
		valid = false
	}

	return valid, model
}
