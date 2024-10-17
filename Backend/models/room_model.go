// Package models provides all the various models for our ORM.
package models

// Room describes our room table and objects.
type Room struct {
	Entity Entity `gorm:"embedded"`
	Parent Parent `gorm:"embedded"`
}
