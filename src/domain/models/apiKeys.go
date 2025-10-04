package models

import (
	"time"

	"github.com/google/uuid"
)

type ApiKeys struct {
	Id          int       `gorm:"type:serial;primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	ProjectID   int       `gorm:"not null" json:"projectId"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Key         string    `gorm:"size:100;not null;unique" json:"key"`
	IpWhiteList []string  `gorm:"type:text" json:"ip_whitelist"`
	ExpireAt    time.Time `gorm:"not null;default:0" json:"expire_at"`
	Revoked     bool      `gorm:"not null;default:false" json:"revoked"`
	RateLimit   int       `gorm:"not null;default:0" json:"rate_limit_per_hour"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
