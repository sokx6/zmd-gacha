package models

type Character struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"not null;uniqueIndex;size:40"`
	Rank string `gorm:"not null;"`
}
