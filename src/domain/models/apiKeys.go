package models

import (
	"github.com/google/uuid"
	"time"
)

type ApiKeys struct {
	Id          int       `gorm:"type:serial;primary_key" json:"id"`
	UserID      uuid.UUID `gorm:"type:serial;not null;unique" json:"userId"`
	Name        string    `gorm:"size:100;not null;unique" json:"name"`
	Key         string    `gorm:"size:100;not null;unique" json:"key"`
	IpWhiteList []string  `gorm:"size:100;not null;unique" json:"ip_whitelist"`
	ExpireAt    time.Time `gorm:"not null;default:0" json:"expire_at"`
	Revoked     bool      `gorm:"not null;default:false" json:"revoked"`
	RateLimit   int       `gorm:"not null;default:0" json:"rate_limit_per_hour"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}
