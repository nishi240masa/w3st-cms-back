package models

import (
	"github.com/google/uuid"
)

type UUID = uuid.UUID

// StringToUUID
func StringToUUID(s string) (UUID, error) {
	return uuid.Parse(s)
}

type User struct {
	ID        UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);not null;unique" json:"email"`
	Password string `gorm:"type:varchar(100);not null" json:"password"`
	CreateAt  string `gorm:"type:timestamp;not null" json:"createAt"`
	UpdateAt  string `gorm:"type:timestamp;not null" json:"updateAt"`
}