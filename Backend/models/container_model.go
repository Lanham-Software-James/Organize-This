// Package models provides all the various models for our ORM.
package models

// Container describes our room table and objects.
type Container struct {
	Entity Entity `gorm:"embedded"`
	Parent Parent `gorm:"embedded"`
}
