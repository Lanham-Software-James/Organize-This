// Package models provides all the various models for our ORM.
package models

// Shelf describes our shelf table and objects.
type Shelf struct {
	Entity Entity `gorm:"embedded"`
}
