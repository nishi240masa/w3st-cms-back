package models

import "time"

type ContentVersion struct {
	ID        UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ContentID UUID      `gorm:"type:uuid;not null" json:"content_id"`
	Version   int       `gorm:"not null" json:"version"`
	Data      string    `gorm:"type:text;not null" json:"data"`
	UserID    UUID      `gorm:"type:uuid;not null" json:"user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
