package repository

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"organize-this/infra/cache"
	"organize-this/infra/database"
	"organize-this/infra/logger"
	"organize-this/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

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

// GetAllEntities returns all entities that belong to the user.
func (repo Repository) GetAllEntities(offset int, limit int) []models.GetEntitiesResponseData {
	stringOffset := strconv.Itoa(offset)
	stringLimit := strconv.Itoa(limit)
	var results []models.GetEntitiesResponseData

	cacheTTL := 5 * time.Minute
	ctx := context.Background()
	key := base64.StdEncoding.EncodeToString([]byte("allentities" + stringOffset + stringLimit))
	value, redisErr := cache.Client.Get(ctx, key).Result()
	if redisErr != nil && !errors.Is(redisErr, redis.Nil) {
		logger.Errorf("Error retriving entites from Redis: %v", redisErr)
		return results
	}

	if value == "" {
		dbErr := repo.Database.Raw(`
			SELECT 'building' AS category, id, name, notes, ' ' as location FROM buildings
			UNION ALL
			SELECT 'room' AS category, id, name, notes, ' ' as location FROM rooms
			UNION ALL
			SELECT 'shelving_unit' AS category, id, name, notes, ' ' as location FROM shelving_units
			UNION ALL
			SELECT 'shelf' AS category, id, name, notes, ' ' as location FROM shelves
			UNION ALL
			SELECT 'container' AS category, id, name, notes, ' ' as location FROM containers
			UNION ALL
			SELECT 'item' AS category, id, name, notes, ' ' as location FROM items
			OFFSET ` + stringOffset +
			` LIMIT ` + stringLimit).Scan(&results).Error

		if dbErr != nil {
			logger.Errorf("error executing query: %v", dbErr)
			return results
		}

		byteData, jsonErr := json.Marshal(results)
		if jsonErr != nil {
			logger.Errorf("error executing query: %v", dbErr)
			return results
		}

		cache.Client.Set(ctx, key, byteData, cacheTTL)
	} else {
		jsonErr := json.Unmarshal([]byte(value), &results)
		if jsonErr != nil {
			logger.Errorf("error executing query: %v", jsonErr)
		}
	}

	return results
}

// CountEntities is used to count the total number of entities that belong to a user.
func (repo Repository) CountEntities() int {
	var entityCount int

	err := repo.Database.Raw(`
		SELECT
			(SELECT COUNT(*) FROM buildings) +
			(SELECT COUNT(*) FROM rooms) +
			(SELECT COUNT(*) FROM shelving_units) +
			(SELECT COUNT(*) FROM shelves) +
			(SELECT COUNT(*) FROM containers) +
			(SELECT COUNT(*) FROM items)
		AS EntityCount
	`).Scan(&entityCount).Error

	if err != nil {
		logger.Errorf("error executing query: %v", err)
	}

	return entityCount
}
