// Package models provides all the various models for our ORM.
package models

import (
	"time"
)

// Shelf describes our shelf table and objects.
type Shelf struct {
	ID        uint64
	Name      string
	Notes     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
