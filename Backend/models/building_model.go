// Package models provides all the various models for our ORM.
package models

// Building describes our building table and objects.
type Building struct {
	Entity  Entity `gorm:"embedded"`
	Address *string
}
