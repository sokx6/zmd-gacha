package models

type Character struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Name      string `gorm:"not null;uniqueIndex,size:40" json:"name"`
	Rank      string `gorm:"not null;" json:"rank"`
	IsLimited bool   `gorm:"not null;default:false;index" json:"is_limited"`
	IsUp      bool   `gorm:"not null;default:false;index" json:"is_up"`
}
