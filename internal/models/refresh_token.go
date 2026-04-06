package models

import "time"

type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	UID       uint   `gorm:"not null"`
	Token     string `gorm:"not null"`
	ExpiredAt time.Time
	Expired   bool
}
