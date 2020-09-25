package model

import (
	"time"

	"gorm.io/gorm"
)

// UserTimezone
type UserTimezone struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ID        string         `gorm:"primaryKey"`
	TimeZone  string
}
