package models

import "time"

type RefreshToken struct {
	UID       uint `gorm:"primaryKey"`
	Token     string
	ExpiredAt time.Time
	Expired   bool
}
