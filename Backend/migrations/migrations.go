package migrations

import (
	"organize-this/infra/database"
	"organize-this/models"
)

func Migrate() {
	var migrationModels = []interface{}{&models.Building{}}
	err := database.GetDB().AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
