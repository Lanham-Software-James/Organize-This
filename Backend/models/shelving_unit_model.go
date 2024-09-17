// Package models provides all the various models for our ORM.
package models

import (
	"time"
)

// ShelvingUnit describes our shelving_unit table and objects.
type ShelvingUnit struct {
	ID        uint64
	Name      string
	Notes     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
