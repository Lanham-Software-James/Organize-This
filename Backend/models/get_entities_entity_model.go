package models

type GetEntitiesEntity struct {
	ID       uint
	Name     string
	Category string
	Parent   []GetEntitiesParentData
	Notes    *string
	Address  *string
}
