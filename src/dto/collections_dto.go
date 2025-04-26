package dto

type MakeCollection struct {
	Name        string `json:"name" binding:"required,min=1"`
	Description string `json:"description" binding:"required,min=1"`
}
