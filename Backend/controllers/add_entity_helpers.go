package controllers

import (
	"organize-this/models"
)

// addItem is a helper function for the CreateEntity endpoint to actually create the item
func (h *Handler) addItem(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	item := models.Item{
		Entity: entity,
	}

	h.Repository.Database.Save(&item)

	return uint(item.Entity.ID)
}

// addContainer is a helper function for the CreateEntity endpoint to actually create the container.
func (h *Handler) addContainer(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	container := models.Container{
		Entity: entity,
	}

	h.Repository.Database.Save(&container)

	return uint(container.Entity.ID)
}

// addShelf is a helper function for the CreateEntity endpoint to actually create the shelf.
func (h *Handler) addShelf(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	shelf := models.Shelf{
		Entity: entity,
	}

	h.Repository.Database.Save(&shelf)

	return uint(shelf.Entity.ID)
}

// addShelvingUnit is a helper function for the CreateEntity endpoint to actually create the shelving unit.
func (h *Handler) addShelvingUnit(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	unit := models.ShelvingUnit{
		Entity: entity,
	}

	h.Repository.Database.Save(&unit)

	return uint(unit.Entity.ID)
}

// addRoom is a helper function for the CreateEntity endpoint to actually create the room.
func (h *Handler) addRoom(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	entity := models.Entity{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	room := models.Room{
		Entity: entity,
	}

	h.Repository.Database.Save(&room)

	return uint(room.Entity.ID)
}

// addBuilding is a helper function for the CreateEntity endpoint to actually create the building.
func (h *Handler) addBuilding(data map[string]string) (id uint) {
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

	h.Repository.Database.Save(&building)

	return uint(building.Entity.ID)
}
