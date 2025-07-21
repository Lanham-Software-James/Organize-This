// Package migrations is used to handle all of our GORM DB migrations.
package migrations

import (
	"willowsuite-vault/infra/database"
	"willowsuite-vault/models"
)

// Migrate is called in main.go to migrate out DB to the latest version.
func Migrate() {
	var migrationModels = []interface{}{&models.Building{}, &models.Room{}, &models.ShelvingUnit{}, &models.Shelf{}, &models.Container{}, &models.Item{}}
	err := database.GetDB().AutoMigrate(migrationModels...)
	if err != nil {
		return
	}
}
