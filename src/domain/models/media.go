package models

import "time"

type MediaAsset struct {
	ID        UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"`
	Path      string    `gorm:"type:text;not null" json:"path"`
	Size      int64     `gorm:"not null" json:"size"`
	UserID    UUID      `gorm:"type:uuid;not null" json:"user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
