package controllers

import (
	"organize-this/models"
)

// addItem is a helper function for the CreateEntity endpoint to actually create the item
func (h *Handler) addItem(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	item := models.Item{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	h.Repository.Database.Save(&item)

	return uint(item.ID)
}

// addContainer is a helper function for the CreateEntity endpoint to actually create the container.
func (h *Handler) addContainer(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	container := models.Container{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	h.Repository.Database.Save(&container)

	return uint(container.ID)
}

// addShelf is a helper function for the CreateEntity endpoint to actually create the shelf.
func (h *Handler) addShelf(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	shelf := models.Shelf{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	h.Repository.Database.Save(&shelf)

	return uint(shelf.ID)
}

// addShelvingUnit is a helper function for the CreateEntity endpoint to actually create the shelving unit.
func (h *Handler) addShelvingUnit(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	unit := models.ShelvingUnit{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	h.Repository.Database.Save(&unit)

	return uint(unit.ID)
}

// addRoom is a helper function for the CreateEntity endpoint to actually create the room.
func (h *Handler) addRoom(data map[string]string) (id uint) {
	tmpNotes := data["notes"]

	room := models.Room{
		Name:  data["name"],
		Notes: &tmpNotes,
	}

	h.Repository.Database.Save(&room)

	return uint(room.ID)
}

// addBuilding is a helper function for the CreateEntity endpoint to actually create the building.
func (h *Handler) addBuilding(data map[string]string) (id uint) {
	tmpAddress := data["address"]
	tmpNotes := data["notes"]

	building := models.Building{
		Name:    data["name"],
		Address: &tmpAddress,
		Notes:   &tmpNotes,
	}

	h.Repository.Database.Save(&building)

	return uint(building.ID)
}
