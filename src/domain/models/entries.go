package models

import (
	"time"
)

type Entry struct {
	ID           int       `gorm:"type:serial;primary_key" json:"id"`
	ProjectID    int       `gorm:"type:int;not null" json:"project_id"`
	CollectionID int       `gorm:"type:int;not null" json:"collection_id"`
	Data         string    `gorm:"type:jsonb" json:"data"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
