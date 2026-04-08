package models

type User struct {
	UID            uint            `gorm:"primarykey"`
	Username       string          `gorm:"not null;uniqueIndex;size:40"`
	Nickname       string          `gorm:"size:40"`
	Profile        string          `gorm:"size:255"`
	Password       string          `gorm:"not null;size:255" json:"password"`
	Email          string          `gorm:"not null;uniqueIndex;size:100"`
	Role           string          `gorm:"not null;default:user;size:20;index"`
	UserCharacters []UserCharacter `gorm:"foreignKey:UserID;references:UID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	GachaRecords   []GachaRecord   `gorm:"foreignKey:UserID;references:UID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PullCount      int             `gorm:"not null;default:0"`
	LastACount     int             `gorm:"not null;default:0"`
	LastSCount     int             `gorm:"not null;default:0"`
	LastSUp        bool            `gorm:"not null;default:false"`
	LastUpCount    int             `gorm:"not null;default:0"`
}
