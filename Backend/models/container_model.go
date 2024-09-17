// Package models provides all the various models for our ORM.
package models

import (
	"time"
)

// Container describes our room table and objects.
type Container struct {
	ID        uint64
	Name      string
	Notes     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
