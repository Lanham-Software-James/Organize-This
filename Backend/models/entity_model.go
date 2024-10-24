// Package models provides all the various models for our ORM.
package models

import (
	"time"

	"gorm.io/gorm"
)

// Entity describes the attributes all entities have in common
type Entity struct {
	ID        uint64
	Name      string
	Notes     *string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
