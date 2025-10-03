package dto

type CreateVersion struct {
	ContentID string `json:"content_id" binding:"required,uuid"`
	Data      string `json:"data" binding:"required"`
	UserID    string `json:"user_id" binding:"required,uuid"`
}

type UpdateVersion struct {
	ID     string `json:"id" binding:"required,uuid"`
	Data   string `json:"data" binding:"omitempty"`
	UserID string `json:"user_id" binding:"omitempty,uuid"`
}

type VersionResponse struct {
	ID        string `json:"id"`
	ContentID string `json:"content_id"`
	Version   int    `json:"version"`
	Data      string `json:"data"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
