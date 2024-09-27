package controllers

import (
	"fmt"
	"organize-this/models"
)

// createEntityByCategory handles the creation of entities based on their category.
func (handler Handler) createEntityByCategory(category string, data map[string]string) (uint, error) {
	switch category {
	case "item":
		return handler.addItem(data), nil
	case "container":
		return handler.addContainer(data), nil
	case "shelf":
		return handler.addShelf(data), nil
	case "shelvingunit":
		return handler.addShelvingUnit(data), nil
	case "room":
		return handler.addRoom(data), nil
	case "building":
		return handler.addBuilding(data), nil
	default:
		return 0, fmt.Errorf("invalid category: %v", category)
	}
}

// addItem is a helper function for the CreateEntity endpoint to actually create the item
func (handler Handler) addItem(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	item := models.Item{
		Entity: entity,
	}

	handler.Repository.Database.Save(&item)

	return uint(item.Entity.ID)
}

// addContainer is a helper function for the CreateEntity endpoint to actually create the container.
func (handler Handler) addContainer(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	container := models.Container{
		Entity: entity,
	}

	handler.Repository.Database.Save(&container)

	return uint(container.Entity.ID)
}

// addShelf is a helper function for the CreateEntity endpoint to actually create the shelf.
func (handler Handler) addShelf(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	shelf := models.Shelf{
		Entity: entity,
	}

	handler.Repository.Database.Save(&shelf)

	return uint(shelf.Entity.ID)
}

// addShelvingUnit is a helper function for the CreateEntity endpoint to actually create the shelving unit.
func (handler Handler) addShelvingUnit(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	unit := models.ShelvingUnit{
		Entity: entity,
	}

	handler.Repository.Database.Save(&unit)

	return uint(unit.Entity.ID)
}

// addRoom is a helper function for the CreateEntity endpoint to actually create the room.
func (handler Handler) addRoom(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	room := models.Room{
		Entity: entity,
	}

	handler.Repository.Database.Save(&room)

	return uint(room.Entity.ID)
}

// addBuilding is a helper function for the CreateEntity endpoint to actually create the building.
func (handler Handler) addBuilding(data map[string]string) (id uint) {
	tmpAddress := data["address"]
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	building := models.Building{
		Entity:  entity,
		Address: &tmpAddress,
	}

	handler.Repository.Database.Save(&building)

	return uint(building.Entity.ID)
}
