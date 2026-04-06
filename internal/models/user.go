package models

type User struct {
	UID        uint        `gorm:"primarykey"`
	Username   string      `gorm:"not null;uniqueIndex;size:40"`
	Password   string      `gorm:"not null;size:40" json:"password"`
	Email      string      `gorm:"not null;uniqueIndex"`
	Characters []Character `gorm:"many2many:user_characters;"`
}
