package repository

import (
	"organize-this/infra/database"
	"organize-this/infra/logger"

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

func (repo Repository) Get(model interface{}, offset int, limit int) interface{} {
	err := repo.Database.Offset(offset).Limit(limit).Find(model).Error
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

func (repo Repository) Count(model interface{}) (int, error) {
	var count int64
	err := repo.Database.Model(model).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
