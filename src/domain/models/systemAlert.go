package models

import (
	"time"
)

type SystemAlert struct {
	ID        int                    `json:"id" gorm:"primaryKey"`
	AlertType string                 `json:"alert_type" gorm:"not null"`
	Severity  string                 `json:"severity" gorm:"not null"`
	Title     string                 `json:"title" gorm:"not null"`
	Message   string                 `json:"message" gorm:"not null"`
	ProjectID int                    `json:"project_id" gorm:"not null;default:1"`
	IsActive  bool                   `json:"is_active" gorm:"not null;default:true"`
	IsRead    bool                   `json:"is_read" gorm:"not null;default:false"`
	Metadata  map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
