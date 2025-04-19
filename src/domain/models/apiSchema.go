package models

type ApiSchema struct {
	ID        int    `gorm:"type:serial;primary_key" json:"id"`
	UserID    UUID   `gorm:"type:uuid;not null" json:"userId"`
	FieldID   string `gorm:"type:varchar(100);not null" json:"fieldId"`
	ViewName  string `gorm:"type:varchar(100);not null" json:"viewName"`
	FieldType string `gorm:"type:varchar(50);not null" json:"fieldType"`
	CreateAt  string `gorm:"type:timestamp;not null" json:"createAt"`
	UpdateAt  string `gorm:"type:timestamp;not null" json:"updateAt"`
}
