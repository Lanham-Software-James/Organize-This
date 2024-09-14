package migrations

import (
	"chi-boilerplate/infra/database"
	"chi-boilerplate/models"
)

func Migrate() {
	var migrationModels = []interface{}{&models.Building{}}
	err := database.GetDB().AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
