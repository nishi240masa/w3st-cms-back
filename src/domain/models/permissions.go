package models

import "time"

type UserPermission struct {
	ID         UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID     UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Permission string    `gorm:"type:varchar(50);not null" json:"permission"`
	Resource   string    `gorm:"type:varchar(255);not null" json:"resource"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}