// Package models provides all the various models for our ORM.
package models

import (
	"time"
)

// Building describes our building table and objects.
type Building struct {
	ID        uint64
	Name      string
	Address   *string
	Notes     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
