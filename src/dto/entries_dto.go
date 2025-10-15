package dto

type CreateEntry struct {
	Data map[string]interface{} `json:"data" binding:"required"`
}

type UpdateEntry struct {
	Data map[string]interface{} `json:"data" binding:"required"`
}
