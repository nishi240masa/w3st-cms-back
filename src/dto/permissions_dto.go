package dto

type CreatePermission struct {
	UserID     string `json:"user_id" binding:"required,uuid"`
	Permission string `json:"permission" binding:"required,min=1"`
	Resource   string `json:"resource" binding:"required,min=1"`
}

type UpdatePermission struct {
	ID         string `json:"id" binding:"required,uuid"`
	Permission string `json:"permission" binding:"omitempty,min=1"`
	Resource   string `json:"resource" binding:"omitempty,min=1"`
}

type PermissionResponse struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Permission string `json:"permission"`
	Resource   string `json:"resource"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
