package models

import "time"

type GachaPool struct {
	ID          uint            `gorm:"primarykey"`
	Name        string          `gorm:"not null;uniqueIndex;size:100"`
	Description string          `gorm:"size:255"`
	Config      GachaPoolConfig `gorm:"foreignKey:PoolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StartAt     *time.Time      `gorm:"index"`
	EndAt       *time.Time      `gorm:"index"`
	IsActive    bool            `gorm:"not null;default:true;index"`
	Characters  []Character     `gorm:"many2many:gacha_pool_characters;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type GachaPoolConfig struct {
	ID                 uint    `gorm:"primarykey"`
	PoolID             uint    `gorm:"not null;uniqueIndex"`
	SRankBaseRate      float64 `gorm:"not null;default:0.008"`
	ARankBaseRate      float64 `gorm:"not null;default:0.08"`
	BRankBaseRate      float64 `gorm:"not null;default:0.912"`
	AGuaranteeInterval int     `gorm:"not null;default:10"`
	SPityStart         int     `gorm:"not null;default:65"`
	SPityStep          float64 `gorm:"not null;default:0.05"`
	SPityEnd           int     `gorm:"not null;default:80"`
	LimitPity          int     `gorm:"not null;default:120"`
	LimitRateWhenS     float64 `gorm:"not null;default:0.5"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type GachaRecord struct {
	ID          uint `gorm:"primarykey"`
	UserID      uint `gorm:"not null;index:idx_user_record"`
	PoolID      uint `gorm:"not null;index:idx_pool_record"`
	CharacterID uint `gorm:"not null;index:idx_char_record"`
	PullCount   int  `gorm:"not null;default:0"`
	CreatedAt   time.Time
}

type UserCharacter struct {
	ID          uint      `gorm:"primarykey"`
	UserID      uint      `gorm:"not null;index:idx_user_char,unique"`
	User        User      `gorm:"foreignKey:UserID;references:UID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CharacterID uint      `gorm:"not null;index:idx_user_char,unique"`
	Character   Character `gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OwnedCount  int       `gorm:"not null;default:0"`
	Level       int       `gorm:"not null;default:0"`
}
