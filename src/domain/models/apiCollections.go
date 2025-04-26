package models

import "time"

type ApiCollection struct {
	ID          int       `gorm:"type:serial;primary_key" json:"id"`
	UserID      UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
