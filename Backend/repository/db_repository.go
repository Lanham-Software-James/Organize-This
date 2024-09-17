package repository

import (
	"chi-boilerplate/infra/database"
	"chi-boilerplate/infra/logger"

	"gorm.io/gorm"
)

type Repository struct {
	Database *gorm.DB
}

func (repo Repository) Save(model interface{}) interface{} {
	err := repo.Database.Create(model).Error

	if err != nil {
		logger.Errorf("error, not save data %v", err)
	}

	return err
}

func (repo Repository) Get(model interface{}) interface{} {
	err := database.DB.Find(model).Error
	return err
}

func (repo Repository) GetOne(model interface{}) interface{} {
	err := database.DB.Last(model).Error
	return err
}

func (repo Repository) Update(model interface{}) interface{} {
	err := database.DB.Find(model).Error
	return err
}
