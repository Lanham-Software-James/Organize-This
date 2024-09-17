package migrations

import (
	"organize-this/infra/database"
	"organize-this/models"
)

func Migrate() {
	var migrationModels = []interface{}{&models.Building{}, &models.Room{}, &models.ShelvingUnit{}, &models.Shelf{}}
	err := database.GetDB().AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
