package dto

type CreateAuditLog struct {
	UserID   string `json:"user_id" binding:"required,uuid"`
	Action   string `json:"action" binding:"required,min=1"`
	Resource string `json:"resource" binding:"required,min=1"`
	Details  string `json:"details"`
}

type AuditLogResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	CreatedAt string `json:"created_at"`
	Details   string `json:"details"`
}
