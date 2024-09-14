package models

import (
	"time"
)

type Building struct {
	ID        uint64
	Name      string
	Address   *string
	Notes     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
