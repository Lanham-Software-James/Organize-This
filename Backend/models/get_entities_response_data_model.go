package models

type GetEntitiesResponseData struct {
	ID             uint
	Name           string
	Category       string
	ParentID       uint
	ParentCategory string
	Notes          *string
}
