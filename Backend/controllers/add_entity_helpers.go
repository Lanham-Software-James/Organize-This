package controllers

import (
	"context"
	"fmt"
	"organize-this/models"
)

// createEntityByCategory handles the creation of entities based on their category.
func (handler Handler) createEntityByCategory(ctx context.Context, userID string, category string, data map[string]string) (uint, error) {
	var id uint

	switch category {
	case "item":
		id = handler.addItem(userID, data)
	case "container":
		id = handler.addContainer(userID, data)
	case "shelf":
		id = handler.addShelf(userID, data)
	case "shelvingunit":
		id = handler.addShelvingUnit(userID, data)
	case "room":
		id = handler.addRoom(userID, data)
	case "building":
		id = handler.addBuilding(userID, data)
	default:
		return 0, fmt.Errorf("invalid category: %v", category)
	}

	handler.Repository.FlushEntities(ctx)
	return id, nil
}

// addItem is a helper function for the CreateEntity endpoint to actually create the item
func (handler Handler) addItem(userID string, data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:   data["name"],
		Notes:  &tmpNotes,
		UserID: userID,
	}

	item := models.Item{
		Entity: entity,
	}

	handler.Repository.Save(&item)

	return uint(item.Entity.ID)
}

// addContainer is a helper function for the CreateEntity endpoint to actually create the container.
func (handler Handler) addContainer(userID string, data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:   data["name"],
		Notes:  &tmpNotes,
		UserID: userID,
	}

	container := models.Container{
		Entity: entity,
	}

	handler.Repository.Save(&container)

	return uint(container.Entity.ID)
}

// addShelf is a helper function for the CreateEntity endpoint to actually create the shelf.
func (handler Handler) addShelf(userID string, data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:   data["name"],
		Notes:  &tmpNotes,
		UserID: userID,
	}

	shelf := models.Shelf{
		Entity: entity,
	}

	handler.Repository.Save(&shelf)

	return uint(shelf.Entity.ID)
}

// addShelvingUnit is a helper function for the CreateEntity endpoint to actually create the shelving unit.
func (handler Handler) addShelvingUnit(userID string, data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:   data["name"],
		Notes:  &tmpNotes,
		UserID: userID,
	}

	unit := models.ShelvingUnit{
		Entity: entity,
	}

	handler.Repository.Save(&unit)

	return uint(unit.Entity.ID)
}

// addRoom is a helper function for the CreateEntity endpoint to actually create the room.
func (handler Handler) addRoom(userID string, data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:   data["name"],
		Notes:  &tmpNotes,
		UserID: userID,
	}

	room := models.Room{
		Entity: entity,
	}

	handler.Repository.Save(&room)

	return uint(room.Entity.ID)
}

// addBuilding is a helper function for the CreateEntity endpoint to actually create the building.
func (handler Handler) addBuilding(userID string, data map[string]string) (id uint) {
	tmpAddress := data["address"]
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:   data["name"],
		Notes:  &tmpNotes,
		UserID: userID,
	}

	building := models.Building{
		Entity:  entity,
		Address: &tmpAddress,
	}

	handler.Repository.Save(&building)

	return uint(building.Entity.ID)
}
