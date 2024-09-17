package models

import (
	"time"
)

type ShelvingUnit struct {
	ID        uint64
	Name      string
	Notes     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
