package dto

type CreateMedia struct {
	Name   string `json:"name" binding:"required,min=1"`
	Type   string `json:"type" binding:"required,min=1"`
	Path   string `json:"path" binding:"required,min=1"`
	Size   int64  `json:"size" binding:"required,min=1"`
	UserID string `json:"user_id" binding:"required,uuid"`
}

type UpdateMedia struct {
	ID   string `json:"id" binding:"required,uuid"`
	Name string `json:"name" binding:"omitempty,min=1"`
	Type string `json:"type" binding:"omitempty,min=1"`
	Path string `json:"path" binding:"omitempty,min=1"`
	Size int64  `json:"size" binding:"omitempty,min=1"`
}

type MediaResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
