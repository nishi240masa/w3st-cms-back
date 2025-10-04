package models

import "time"

type FieldData struct {
	ID           int       `gorm:"type:serial;primary_key" json:"id"`
	ProjectID    int       `gorm:"type:int;not null" json:"project_id"`
	CollectionID int       `gorm:"type:int;not null" json:"collection_id"`
	FieldID      string    `gorm:"type:varchar(100);not null" json:"field_id"`
	ViewName     string    `gorm:"type:varchar(100);not null" json:"view_name"`
	FieldType    string    `gorm:"type:varchar(50);not null" json:"field_type"`
	IsRequired   bool      `gorm:"not null;default:false" json:"is_required"`
	DefaultValue string    `gorm:"type:jsonb" json:"default_value"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
