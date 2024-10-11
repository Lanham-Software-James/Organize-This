// Package repository contains all functions for data manipulation for Organize-This
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"organize-this/infra/cache"
	"organize-this/infra/logger"
	"organize-this/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

// Repository is the type we use to pass the infrastructure components to the functions
type Repository struct {
	Database *gorm.DB
	Cache    *redis.Client
}

// GetEntitiesCacheKey is an extension of cachekey that represents the structure of the keys in our cache for the paginated getentities data.
type GetEntitiesCacheKey struct {
	CacheKey cache.CacheKey
	Offset   string
	Limit    string
}

// Save is used to create a new record in the DB
func (repo Repository) Save(model interface{}) interface{} {
	err := repo.Database.Save(model).Error

	if err != nil {
		logger.Errorf("error, not save data %v", err)
	}

	return err
}

// GetOne is used to get a single record from the DB
func (repo Repository) GetOne(model interface{}, userID string) interface{} {
	err := repo.Database.Where("user_id = ?", userID).First(model).Error
	return err
}

// GetAllEntities returns all entities that belong to the user.
func (repo Repository) GetAllEntities(ctx context.Context, userID string, offset int, limit int) []models.GetEntitiesResponseData {
	stringOffset := strconv.Itoa(offset)
	stringLimit := strconv.Itoa(limit)
	var results []models.GetEntitiesResponseData

	cacheTTL := 5 * time.Minute
	keyStructured := GetEntitiesCacheKey{
		CacheKey: cache.CacheKey{
			User:     userID,
			Function: "GetAllEntities",
		},
		Offset: stringOffset,
		Limit:  stringLimit,
	}
	key, _ := json.Marshal(keyStructured)
	value, redisErr := repo.Cache.Get(ctx, string(key)).Result()
	if redisErr != nil && !errors.Is(redisErr, redis.Nil) {
		logger.Errorf("Error retriving entites from Redis: %v", redisErr)
		return results
	}

	if value == "" {
		dbErr := repo.Database.Raw(`
			(SELECT 'building' AS category, id, name, notes, ' ' as location FROM buildings WHERE user_id = ? LIMIT `+stringLimit+`)
			UNION ALL
			(SELECT 'room' AS category, id, name, notes, ' ' as location FROM rooms WHERE user_id = ? LIMIT `+stringLimit+`)
			UNION ALL
			(SELECT 'shelving_unit' AS category, id, name, notes, ' ' as location FROM shelving_units WHERE user_id = ? LIMIT `+stringLimit+`)
			UNION ALL
			(SELECT 'shelf' AS category, id, name, notes, ' ' as location FROM shelves WHERE user_id = ? LIMIT `+stringLimit+`)
			UNION ALL
			(SELECT 'container' AS category, id, name, notes, ' ' as location FROM containers WHERE user_id = ? LIMIT `+stringLimit+`)
			UNION ALL
			(SELECT 'item' AS category, id, name, notes, ' ' as location FROM items WHERE user_id = ? LIMIT `+stringLimit+`)
			 OFFSET `+stringOffset+
			` LIMIT `+stringLimit, userID, userID, userID, userID, userID, userID).Scan(&results).Error

		if dbErr != nil {
			logger.Errorf("error executing query: %v", dbErr)
			return results
		}

		byteData, jsonErr := json.Marshal(results)
		if jsonErr != nil {
			logger.Errorf("error executing query: %v", dbErr)
			return results
		}

		repo.Cache.Set(ctx, string(key), byteData, cacheTTL)
	} else {
		jsonErr := json.Unmarshal([]byte(value), &results)
		if jsonErr != nil {
			logger.Errorf("error executing query: %v", jsonErr)
		}
	}

	return results
}

// CountEntities is used to count the total number of entities that belong to a user.
func (repo Repository) CountEntities(ctx context.Context, userID string) int {
	var entityCount int

	cacheTTL := 5 * time.Minute
	keyStructured := cache.CacheKey{
		User:     userID,
		Function: "CountEntities",
	}
	key, _ := json.Marshal(keyStructured)
	value, redisErr := repo.Cache.Get(ctx, string(key)).Result()
	if redisErr != nil && !errors.Is(redisErr, redis.Nil) {
		logger.Errorf("error retriving entites from Redis: %v", redisErr)
		return entityCount
	}

	if value == "" {
		err := repo.Database.Raw(`
			SELECT
				(SELECT COUNT(*) FROM buildings WHERE user_id = ?) +
				(SELECT COUNT(*) FROM rooms WHERE user_id = ?) +
				(SELECT COUNT(*) FROM shelving_units WHERE user_id = ?) +
				(SELECT COUNT(*) FROM shelves WHERE user_id = ?) +
				(SELECT COUNT(*) FROM containers WHERE user_id = ?) +
				(SELECT COUNT(*) FROM items WHERE user_id = ?)
			AS EntityCount
		`, userID, userID, userID, userID, userID, userID).Scan(&entityCount).Error

		if err != nil {
			logger.Errorf("error executing query: %v", err)
			return entityCount
		}

		repo.Cache.Set(ctx, string(key), entityCount, cacheTTL)
	} else {
		var typeErr error
		entityCount, typeErr = strconv.Atoi(value)
		if typeErr != nil {
			logger.Errorf("error converting string to int from cache: %v", typeErr)
			return entityCount
		}
	}

	return entityCount
}

// FlushEntities clears the redis cache of all things relating to entities
func (repo Repository) FlushEntities(ctx context.Context, userID string) {

	getAllEntitiesPattern := `{"CacheKey":{"User":"` + userID + `","Function":"GetAllEntities"},*`
	countEntitiesPattern := `{"User":"` + userID + `","Function":"CountEntities"}`

	keys, err := repo.Cache.Keys(ctx, getAllEntitiesPattern).Result()
	if err != nil {
		logger.Errorf("error getting cache keys: %v", err)
		return
	}

	for _, key := range keys {
		err := repo.Cache.Del(ctx, key).Err()
		if err != nil {
			logger.Errorf("error clearing cache: %v", err)
			return
		}
	}

	err = repo.Cache.Del(ctx, countEntitiesPattern).Err()
	if err != nil {
		logger.Errorf("error clearing cache: %v", err)
		return
	}
}
