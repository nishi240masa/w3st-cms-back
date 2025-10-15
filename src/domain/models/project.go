package models

import (
	"time"
)

type Project struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	Name             string    `json:"name" gorm:"not null"`
	Description      string    `json:"description"`
	RateLimitPerHour int       `json:"rate_limit_per_hour" gorm:"default:1000"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
