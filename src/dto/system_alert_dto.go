package dto

type CreateSystemAlertRequest struct {
	AlertType string                 `json:"alert_type" binding:"required"`
	Severity  string                 `json:"severity" binding:"required,oneof=info warning error critical"`
	Title     string                 `json:"title" binding:"required,max=255"`
	Message   string                 `json:"message" binding:"required"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type SystemAlertResponse struct {
	ID        int                    `json:"id"`
	AlertType string                 `json:"alert_type"`
	Severity  string                 `json:"severity"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	ProjectID int                    `json:"project_id"`
	IsActive  bool                   `json:"is_active"`
	IsRead    bool                   `json:"is_read"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

type GetSystemAlertsQuery struct {
	Limit  int `form:"limit,default=50"`
	Offset int `form:"offset,default=0"`
}

type SystemAlertCountResponse struct {
	Count int `json:"count"`
}