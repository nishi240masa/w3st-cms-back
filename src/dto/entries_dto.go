package dto

type CreateEntry struct {
	Data map[string]interface{} `json:"data" binding:"required"`
}