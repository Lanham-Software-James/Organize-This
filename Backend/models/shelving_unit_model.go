// Package models provides all the various models for our ORM.
package models

// ShelvingUnit describes our shelving_unit table and objects.
type ShelvingUnit struct {
	Entity Entity `gorm:"embedded"`
}
