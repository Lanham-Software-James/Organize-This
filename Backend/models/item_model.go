// Package models provides all the various models for our ORM.
package models

// Item describes our room table and objects.
type Item struct {
	Entity Entity `gorm:"embedded"`
}
