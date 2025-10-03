package models

import "time"

type AuditLog struct {
	ID        UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Action    string    `gorm:"type:varchar(50);not null" json:"action"`
	Resource  string    `gorm:"type:varchar(255);not null" json:"resource"`
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"timestamp"`
	Details   string    `gorm:"type:text" json:"details"`
}