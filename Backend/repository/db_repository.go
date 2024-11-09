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
	"slices"
	"strconv"
	"strings"
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
	Search   string
	Filters  []string
}

// CountEntitiesCacheKey is an extension of cachekey that represents the structure of the keys in our cache for the count entities data.
type CountEntitiesCacheKey struct {
	CacheKey cache.CacheKey
	Search   string
	Filters  []string
}

type getEntitiesTables struct {
	category    string
	tableName   string
	tableWeight int
}

type countEntitiesTables struct {
	category  string
	tableName string
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
func (repo Repository) GetAllEntities(ctx context.Context, userID string, offset int, limit int, search string, filters []string) ([]models.GetEntitiesEntity, error) {
	stringOffset := strconv.Itoa(offset)
	stringLimit := strconv.Itoa(limit)
	var data []models.GetEntitiesEntity

	cacheTTL := 5 * time.Minute
	keyStructured := GetEntitiesCacheKey{
		CacheKey: cache.CacheKey{
			User:     userID,
			Function: "GetAllEntities",
		},
		Offset:  stringOffset,
		Limit:   stringLimit,
		Search:  search,
		Filters: filters,
	}
	key, jsonErr := json.Marshal(keyStructured)
	if jsonErr != nil {
		logger.Errorf("Error encoding Redis key: %v", jsonErr)
		return nil, jsonErr
	}

	value, redisErr := repo.Cache.Get(ctx, string(key)).Result()
	if redisErr != nil && !errors.Is(redisErr, redis.Nil) {
		logger.Errorf("Error retriving entites from Redis: %v", redisErr)
		return nil, redisErr
	}

	if value == "" {
		var results []models.GetEntitiesResponseData
		stringValues := []string{}
		mainSQL := []string{}

		searchSQL := ""
		buildingSearchSQL := ""
		addSearch := false

		// Dynamically build search query
		if search != "" {
			search = strings.ToLower(search)
			search = "%" + search + "%"
			searchSQL = " AND (LOWER(name) LIKE ? OR LOWER(notes) LIKE ?)"
			buildingSearchSQL = " AND (LOWER(name) LIKE ? OR LOWER(notes) LIKE ? OR LOWER(address) LIKE ?)"
			addSearch = true
		}

		tables := []getEntitiesTables{
			{"building", "buildings", 1},
			{"room", "rooms", 2},
			{"shelving_unit", "shelving_units", 3},
			{"shelf", "shelves", 4},
			{"container", "containers", 5},
			{"item", "items", 6},
		}

		for _, table := range tables {
			if len(filters) == 0 || slices.Contains(filters, table.category) {
				query := ""
				if table.category == "building" {
					query = fmt.Sprintf(`(SELECT %d AS tableWeight, created_at, '%s' AS category, id, name, notes, address, 0 AS parent_id, ' ' AS parent_category FROM %s WHERE user_id = ? AND deleted_at IS NULL`, table.tableWeight, table.category, table.tableName)
				} else {
					query = fmt.Sprintf(`(SELECT %d AS tableWeight, created_at, '%s' AS category, id, name, notes, '' AS address, parent_id, parent_category FROM %s WHERE user_id = ? AND deleted_at IS NULL`, table.tableWeight, table.category, table.tableName)
				}

				stringValues = append(stringValues, userID)

				if addSearch && table.category == "building" {
					query += buildingSearchSQL
					stringValues = append(stringValues, search, search, search)
				} else if addSearch {
					query += searchSQL
					stringValues = append(stringValues, search, search)
				}

				query += ` ORDER BY created_at ASC LIMIT ?)`
				stringValues = append(stringValues, stringLimit)

				mainSQL = append(mainSQL, query)
			}
		}

		// Union all dynamically built queries
		unionQuery := strings.Join(mainSQL, " UNION ALL ")

		// Default Order
		if len(filters) == 0 || len(filters) > 1 {
			unionQuery += " ORDER BY tableWeight"
		}

		// Add offset
		unionQuery += " OFFSET ?"
		stringValues = append(stringValues, stringOffset)

		// Add union limit if more than one filter applied
		if len(filters) == 0 || len(filters) > 1 {
			unionQuery += " LIMIT ?"
			stringValues = append(stringValues, stringLimit)
		}

		// Convert string slice to interface slice
		values := make([]interface{}, len(stringValues))
		for i, stringID := range stringValues {
			values[i] = stringID
		}

		// Run dynamically built query
		dbErr := repo.Database.Raw(unionQuery, values...).Scan(&results).Error
		if dbErr != nil {
			logger.Errorf("error executing query: %v", dbErr)
			return nil, dbErr
		}

		// Generate results
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

		// Set cache
		byteData, jsonErr := json.Marshal(data)
		if jsonErr != nil {
			logger.Errorf("error encoding data: %v", dbErr)
			return nil, jsonErr
		}

		repo.Cache.Set(ctx, string(key), byteData, cacheTTL)
	} else {
		jsonErr := json.Unmarshal([]byte(value), &data)
		if jsonErr != nil {
			logger.Errorf("error encoding data: %v", jsonErr)
		}
	}

	return data, nil
}

// CountEntities is used to count the total number of entities that belong to a user.
func (repo Repository) CountEntities(ctx context.Context, userID string, search string, filters []string) int {
	var entityCount int

	cacheTTL := 5 * time.Minute
	keyStructured := CountEntitiesCacheKey{
		CacheKey: cache.CacheKey{
			User:     userID,
			Function: "CountEntities",
		},
		Search:  search,
		Filters: filters,
	}

	key, _ := json.Marshal(keyStructured)
	value, redisErr := repo.Cache.Get(ctx, string(key)).Result()
	if redisErr != nil && !errors.Is(redisErr, redis.Nil) {
		logger.Errorf("error retriving entites from Redis: %v", redisErr)
		return entityCount
	}

	if value == "" {
		searchSQL := ""
		buildingSearchSQL := ""
		addSearch := false

		stringValues := []string{}
		mainSQL := []string{}

		if search != "" {
			search = strings.ToLower(search)
			search = "%" + search + "%"
			searchSQL = "AND (LOWER(name) LIKE ? OR LOWER(notes) LIKE ?)"
			buildingSearchSQL = "AND (LOWER(name) LIKE ? OR LOWER(notes) LIKE ? OR LOWER(address) LIKE ?)"
			addSearch = true
		}

		tables := []countEntitiesTables{
			{"building", "buildings"},
			{"room", "rooms"},
			{"shelving_unit", "shelving_units"},
			{"shelf", "shelves"},
			{"container", "containers"},
			{"item", "items"},
		}

		for _, table := range tables {
			if len(filters) == 0 || slices.Contains(filters, table.category) {
				query := fmt.Sprintf(`(SELECT COUNT(*) FROM %s WHERE user_id = ? AND deleted_at IS NULL `, table.tableName)

				stringValues = append(stringValues, userID)

				if addSearch && table.category == "building" {
					query += buildingSearchSQL
					stringValues = append(stringValues, search, search, search)
				} else if addSearch {
					query += searchSQL
					stringValues = append(stringValues, search, search)
				}

				query += `)`

				mainSQL = append(mainSQL, query)
			}
		}

		// Union all dynamically built queries
		joinedQueries := strings.Join(mainSQL, " + ")

		additionQuery := "SELECT " + joinedQueries + " AS EntityCount"

		// Convert string slice to interface slice
		values := make([]interface{}, len(stringValues))
		for i, stringID := range stringValues {
			values[i] = stringID
		}

		err := repo.Database.Raw(additionQuery, values...).Scan(&entityCount).Error
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
	countEntitiesPattern := `{"CacheKey":{"User":"` + userID + `","Function":"CountEntities"},*`

	patterns := []string{
		`{"User":"` + userID + `","Function":"GetItemParents"}`,
		`{"User":"` + userID + `","Function":"GetContainerParents"}`,
		`{"User":"` + userID + `","Function":"GetShelfParents"}`,
		`{"User":"` + userID + `","Function":"GetShelving_unitParents"}`,
		`{"User":"` + userID + `","Function":"GetRoomParents"}`,
	}

	entitiesKeys, err := repo.Cache.Keys(ctx, getAllEntitiesPattern).Result()
	if err != nil {
		logger.Errorf("error getting cache keys: %v", err)
		return
	}

	countKeys, err := repo.Cache.Keys(ctx, countEntitiesPattern).Result()
	if err != nil {
		logger.Errorf("error getting cache keys: %v", err)
		return
	}

	keys := append(entitiesKeys, countKeys...)
	keys = append(keys, patterns...)

	for _, key := range keys {
		err := repo.Cache.Del(ctx, key).Err()
		if err != nil {
			logger.Errorf("error clearing cache: %v", err)
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
