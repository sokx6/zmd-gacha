package models

import "time"

type GachaPool struct {
	ID                  uint                 `gorm:"primarykey" json:"id"`
	Name                string               `gorm:"not null;uniqueIndex;size:100" json:"name"`
	Description         string               `gorm:"size:255" json:"description"`
	Config              GachaPoolConfig      `gorm:"foreignKey:PoolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	GachaPoolCharacters []GachaPoolCharacter `gorm:"foreignKey:PoolID;references:ID"`
	Characters          []Character          `gorm:"many2many:gacha_pool_characters;joinForeignKey:PoolID;joinReferences:CharacterID"`
	StartAt             *time.Time           `gorm:"index" json:"start_at"`
	EndAt               *time.Time           `gorm:"index" json:"end_at"`
	IsActive            bool                 `gorm:"not null;default:true;index" json:"is_active"`
}

type GachaPoolCharacter struct {
	ID     uint      `gorm:"primarykey"`
	PoolID uint      `gorm:"not null;uniqueIndex:pool_char"`
	Pool   GachaPool `gorm:"foreignKey:PoolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CharacterID uint      `gorm:"not null;uniqueIndex:pool_char"`
	Character   Character `gorm:"foreignKey:CharacterID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time
}

type GachaPoolConfig struct {
	ID                   uint      `gorm:"primarykey" json:"id"`
	PoolID               uint      `gorm:"not null;uniqueIndex" json:"pool_id"`
	SRankBaseRate        float64   `gorm:"not null;default:0.008" json:"s_rank_base_rate"`
	ARankBaseRate        float64   `gorm:"not null;default:0.08" json:"a_rank_base_rate"`
	AGuaranteeInterval   int       `gorm:"not null;default:10" json:"a_guarantee_interval"`
	SPityStart           int       `gorm:"not null;default:65" json:"s_pity_start"`
	SPityStep            float64   `gorm:"not null;default:0.05" json:"s_pity_step"`
	SPityEnd             int       `gorm:"not null;default:80" json:"s_pity_end"`
	LimitPity            int       `gorm:"not null;default:120" json:"limit_pity"`
	LimitRateWhenS       float64   `gorm:"not null;default:0.5" json:"limit_rate_when_s"`
	MaxLimitedCharacters int       `gorm:"not null;default:0" json:"max_limited_characters"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
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
