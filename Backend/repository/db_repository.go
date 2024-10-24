// Package repository contains all functions for data manipulation for Organize-This
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"organize-this/infra/cache"
	"organize-this/infra/logger"
	"organize-this/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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
func (repo Repository) GetAllEntities(ctx context.Context, userID string, offset int, limit int) ([]models.GetEntitiesEntity, error) {
	stringOffset := strconv.Itoa(offset)
	stringLimit := strconv.Itoa(limit)
	var data []models.GetEntitiesEntity

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
		return nil, redisErr
	}

	if value == "" {
		var results []models.GetEntitiesResponseData
		dbErr := repo.Database.Raw(`
			(SELECT 'building' AS category, id, name, notes, 0 AS parent_id, ' ' AS parent_category FROM buildings WHERE user_id = ? AND deleted_at IS NULL LIMIT `+stringLimit+`)
			UNION ALL
			(SELECT 'room' AS category, id, name, notes, parent_id, parent_category FROM rooms WHERE user_id = ? AND deleted_at IS NULL LIMIT `+stringLimit+`)
			UNION ALL
			(SELECT 'shelving_unit' AS category, id, name, notes, parent_id, parent_category FROM shelving_units WHERE user_id = ? AND deleted_at IS NULL LIMIT `+stringLimit+`)
			UNION ALL
			(SELECT 'shelf' AS category, id, name, notes, parent_id, parent_category FROM shelves WHERE user_id = ? AND deleted_at IS NULL LIMIT `+stringLimit+`)
			UNION ALL
			(SELECT 'container' AS category, id, name, notes, parent_id, parent_category FROM containers WHERE user_id = ? AND deleted_at IS NULL LIMIT `+stringLimit+`)
			UNION ALL
			(SELECT 'item' AS category, id, name, notes, parent_id, parent_category FROM items WHERE user_id = ? AND deleted_at IS NULL LIMIT `+stringLimit+`)
			 OFFSET `+stringOffset+
			` LIMIT `+stringLimit, userID, userID, userID, userID, userID, userID).Scan(&results).Error

		if dbErr != nil {
			logger.Errorf("error executing query: %v", dbErr)
			return nil, dbErr
		}

		for _, entity := range results {
			if entity.ParentID != 0 && entity.ParentCategory != "" {
				var parents []models.GetEntitiesParentData
				repo.getParents(entity.ParentID, entity.ParentCategory, userID, &parents)

				data = append(data, models.GetEntitiesEntity{
					ID:       entity.ID,
					Name:     entity.Name,
					Category: entity.Category,
					Parent:   parents,
					Notes:    entity.Notes,
				})
			} else {
				var parents []models.GetEntitiesParentData

				parents = append(parents, models.GetEntitiesParentData{
					ID:       0,
					Name:     "-",
					Category: "",
				})
				data = append(data, models.GetEntitiesEntity{
					ID:       entity.ID,
					Name:     entity.Name,
					Category: entity.Category,
					Parent:   parents,
					Address:  entity.Address,
					Notes:    entity.Notes,
				})
			}
		}

		byteData, jsonErr := json.Marshal(data)
		if jsonErr != nil {
			logger.Errorf("error executing query: %v", dbErr)
			return nil, jsonErr
		}

		repo.Cache.Set(ctx, string(key), byteData, cacheTTL)
	} else {
		jsonErr := json.Unmarshal([]byte(value), &data)
		if jsonErr != nil {
			logger.Errorf("error executing query: %v", jsonErr)
		}
	}

	return data, nil
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
				(SELECT COUNT(*) FROM buildings WHERE user_id = ? AND deleted_at IS NULL) +
				(SELECT COUNT(*) FROM rooms WHERE user_id = ? AND deleted_at IS NULL) +
				(SELECT COUNT(*) FROM shelving_units WHERE user_id = ? AND deleted_at IS NULL) +
				(SELECT COUNT(*) FROM shelves WHERE user_id = ? AND deleted_at IS NULL) +
				(SELECT COUNT(*) FROM containers WHERE user_id = ? AND deleted_at IS NULL) +
				(SELECT COUNT(*) FROM items WHERE user_id = ? AND deleted_at IS NULL)
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

// Delete is used to soft delete a record from the DB
func (repo Repository) Delete(model interface{}, userID string) interface{} {
	err := repo.Database.Where("user_id = ?", userID).Delete(model).Error
	return err
}

// GetParents returns all the possible parents for an item.
func (repo Repository) GetParents(ctx context.Context, category string, userID string) ([]models.GetEntitiesParentData, error) {
	var results []models.GetEntitiesParentData

	caser := cases.Title(language.AmericanEnglish)
	capitalizedCategory := caser.String(category)

	cacheTTL := 5 * time.Minute
	keyStructured := cache.CacheKey{
		User:     userID,
		Function: "Get" + capitalizedCategory + "Parents",
	}

	key, _ := json.Marshal(keyStructured)
	value, redisErr := repo.Cache.Get(ctx, string(key)).Result()
	if redisErr != nil && !errors.Is(redisErr, redis.Nil) {
		logger.Errorf("Error retriving entites from Redis: %v", redisErr)
		return nil, redisErr
	}

	if value == "" {
		var dbErr error
		switch category {
		case "item":
			dbErr = repo.Database.Raw(`
			(SELECT 'room' AS category, id, name FROM rooms WHERE user_id = ? AND deleted_at IS NULL)
			UNION ALL
			(SELECT 'shelf' AS category, id, name FROM shelves WHERE user_id = ? AND deleted_at IS NULL)
			UNION ALL
			(SELECT 'container' AS category, id, name FROM containers WHERE user_id = ? AND deleted_at IS NULL)`,
				userID, userID, userID).Scan(&results).Error
			break
		case "container":
			dbErr = repo.Database.Raw(`
			(SELECT 'room' AS category, id, name FROM rooms WHERE user_id = ? AND deleted_at IS NULL)
			UNION ALL
			(SELECT 'shelf' AS category, id, name FROM shelves WHERE user_id = ? AND deleted_at IS NULL)`,
				userID, userID).Scan(&results).Error
			break
		case "shelf":
			dbErr = repo.Database.Raw(`
			(SELECT 'shelving_unit' AS category, id, name FROM shelving_units WHERE user_id = ? AND deleted_at IS NULL)`,
				userID).Scan(&results).Error
			break
		case "shelving_unit":
			dbErr = repo.Database.Raw(`
			(SELECT 'room' AS category, id, name FROM rooms WHERE user_id = ? AND deleted_at IS NULL)`,
				userID).Scan(&results).Error
			break
		case "room":
			dbErr = repo.Database.Raw(`
			(SELECT 'building' AS category, id, name FROM buildings WHERE user_id = ? AND deleted_at IS NULL)`,
				userID).Scan(&results).Error
			break
		default:
			logger.Errorf("Invalid category for entity.")
			return nil, fmt.Errorf("Invalid category: %v", category)
		}

		if dbErr != nil {
			logger.Errorf("error executing query: %v", dbErr)
			return nil, dbErr
		}

		byteData, jsonErr := json.Marshal(results)
		if jsonErr != nil {
			logger.Errorf("error executing query: %v", dbErr)
			return nil, jsonErr
		}

		repo.Cache.Set(ctx, string(key), byteData, cacheTTL)
	} else {
		jsonErr := json.Unmarshal([]byte(value), &results)
		if jsonErr != nil {
			logger.Errorf("error executing query: %v", jsonErr)
		}
	}

	return results, nil
}

// HasChildren returns all the children for an entity.
func (repo Repository) HasChildren(id uint64, category string, userID string) (bool, int, error) {
	var childrenCount int
	hasChildren := false

	var dbErr error
	switch category {
	case "building":
		dbErr = repo.Database.Raw(`
			(SELECT count(id) AS childrenCount FROM rooms WHERE user_id = ? AND parent_id = ? AND parent_category = ? AND deleted_at IS NULL)`,
			userID,
			id,
			category,
		).Scan(&childrenCount).Error
		break
	case "room":
		dbErr = repo.Database.Raw(`
			SELECT
				(SELECT count(id) FROM shelving_units WHERE user_id = ? AND parent_id = ? AND parent_category = ? AND deleted_at IS NULL) +
				(SELECT count(id) FROM containers WHERE user_id = ? AND parent_id = ? AND parent_category = ? AND deleted_at IS NULL) +
				(SELECT count(id) FROM items WHERE user_id = ? AND parent_id = ? AND parent_category = ? AND deleted_at IS NULL)
			AS childrenCount`,
			userID, id, category,
			userID, id, category,
			userID, id, category,
		).Scan(&childrenCount).Error
		break
	case "shelving_unit":
		dbErr = repo.Database.Raw(`
			(SELECT count(id) AS childrenCount FROM shelves WHERE user_id = ? AND parent_id = ? AND parent_category = ? AND deleted_at IS NULL)`,
			userID,
			id,
			category,
		).Scan(&childrenCount).Error
		break
	case "shelf":
		dbErr = repo.Database.Raw(`
			SELECT
				(SELECT count(id) FROM containers WHERE user_id = ? AND parent_id = ? AND parent_category = ? AND deleted_at IS NULL) +
				(SELECT count(id) FROM items WHERE user_id = ? AND parent_id = ? AND parent_category = ? AND deleted_at IS NULL)
			AS childrenCount`,
			userID, id, category,
			userID, id, category,
		).Scan(&childrenCount).Error
		break
	case "container":
		dbErr = repo.Database.Raw(`
			(SELECT count(id) AS childrenCount FROM items WHERE user_id = ? AND parent_id = ? AND parent_category = ? AND deleted_at IS NULL)`,
			userID,
			id,
			category,
		).Scan(&childrenCount).Error
		break
	default:
		logger.Errorf("Invalid category for retriving children.")
		return false, 0, fmt.Errorf("Invalid category for retriving children: %v", category)
	}

	if dbErr != nil {
		logger.Errorf("error executing query: %v", dbErr)
		return false, 0, dbErr
	}

	if childrenCount > 0 {
		hasChildren = true
	}

	return hasChildren, childrenCount, nil
}

// FlushEntities clears the redis cache of all things relating to entities
func (repo Repository) FlushEntities(ctx context.Context, userID string) {

	getAllEntitiesPattern := `{"CacheKey":{"User":"` + userID + `","Function":"GetAllEntities"},*`

	patterns := []string{
		`{"User":"` + userID + `","Function":"CountEntities"}`,
		`{"User":"` + userID + `","Function":"GetItemParents"}`,
		`{"User":"` + userID + `","Function":"GetContainerParents"}`,
		`{"User":"` + userID + `","Function":"GetShelfParents"}`,
		`{"User":"` + userID + `","Function":"GetShelving_unitParents"}`,
		`{"User":"` + userID + `","Function":"GetRoomParents"}`,
	}

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

	for name, key := range patterns {
		err = repo.Cache.Del(ctx, key).Err()
		if err != nil {
			logger.Errorf("error clearing cache key: %v. Error: %v", name, err)
			return
		}
	}
}

func (repo Repository) getParents(parentID uint, parentCategory string, userID string, array *[]models.GetEntitiesParentData) error {
	findEntity := models.Entity{
		ID: uint64(parentID),
	}

	var parent models.GetEntitiesParentData
	var recursiveID uint
	var recursiveCategory string
	switch parentCategory {
	case "container":
		model := &models.Container{
			Entity: findEntity,
		}

		err := repo.GetOne(&model, userID)
		if err != nil {
			logger.Errorf("error executing query: %v", err)
			return nil
		}

		parent = models.GetEntitiesParentData{
			ID:       parentID,
			Name:     model.Entity.Name,
			Category: parentCategory,
		}

		recursiveID = uint(model.Parent.ParentID)
		recursiveCategory = model.Parent.ParentCategory
		break
	case "shelf":
		model := &models.Shelf{
			Entity: findEntity,
		}
		err := repo.GetOne(&model, userID)
		if err != nil {
			logger.Errorf("error executing query: %v", err)
			return nil
		}

		parent = models.GetEntitiesParentData{
			ID:       parentID,
			Name:     model.Entity.Name,
			Category: parentCategory,
		}
		recursiveID = uint(model.Parent.ParentID)
		recursiveCategory = model.Parent.ParentCategory
		break
	case "shelving_unit":
		model := &models.ShelvingUnit{
			Entity: findEntity,
		}
		err := repo.GetOne(&model, userID)
		if err != nil {
			logger.Errorf("error executing query: %v", err)
			return nil
		}

		parent = models.GetEntitiesParentData{
			ID:       parentID,
			Name:     model.Entity.Name,
			Category: parentCategory,
		}
		recursiveID = uint(model.Parent.ParentID)
		recursiveCategory = model.Parent.ParentCategory
		break
	case "room":
		model := &models.Room{
			Entity: findEntity,
		}
		err := repo.GetOne(&model, userID)
		if err != nil {
			logger.Errorf("error executing query: %v", err)
			return nil
		}

		parent = models.GetEntitiesParentData{
			ID:       parentID,
			Name:     model.Entity.Name,
			Category: parentCategory,
		}
		recursiveID = uint(model.Parent.ParentID)
		recursiveCategory = model.Parent.ParentCategory
		break
	case "building":
		model := &models.Building{
			Entity: findEntity,
		}
		err := repo.GetOne(&model, userID)
		if err != nil {
			logger.Errorf("error executing query: %v", err)
			return nil
		}

		parent = models.GetEntitiesParentData{
			ID:       parentID,
			Name:     model.Entity.Name,
			Category: parentCategory,
		}
		break
	default:
		logger.Errorf("Invalid Category: %v", parentCategory)
	}

	*array = append((*array), parent)

	if parent.Category != "building" && recursiveID != 0 && recursiveCategory != "" {
		repo.getParents(recursiveID, recursiveCategory, userID, array)
	}

	return nil
}
