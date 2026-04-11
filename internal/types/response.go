package types

import (
	"time"
	"zmd-gacha/internal/models"
)

type UserRstRsp struct {
	Code    int    `json:"code,omitempty"`
	UID     uint   `json:"uid"`
	Message string `json:"message"`
	Role    string `json:"role"`
}

type UserLoginRsp struct {
	Code         int    `json:"code,omitempty"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type TokenRefRsp struct {
	Code         int    `json:"code,omitempty"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type ProfileUpdateRsp struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

type ErrorRsp struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

type PullOnceRsp struct {
	models.Character
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

type PullTenRsp struct {
	Characters []models.Character `json:"characters"`
	Code       int                `json:"code,omitempty"`
	Message    string             `json:"message"`
}

type CharCreateRsp struct {
	models.Character
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

type PoolCreateRsp struct {
	Code    int    `json:"code,omitempty"`
	PoolID  uint   `json:"pool_id"`
	Message string `json:"message"`
}

type InsertCharRsp struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

type CharsGetRsp struct {
	Characters []UserCharacterRsp `json:"characters"`
	Code       int                `json:"code,omitempty"`
	Message    string             `json:"message"`
}

type PoolInfoRsp struct {
	Pool    GachaPoolRsp `json:"pool"`
	Code    int          `json:"code,omitempty"`
	Message string       `json:"message"`
}

type UserRsp struct {
	UID            uint          `json:"UID"`
	Username       string        `json:"Username"`
	Nickname       string        `json:"Nickname"`
	Profile        string        `json:"Profile"`
	Password       string        `json:"password"`
	Email          string        `json:"Email"`
	Role           string        `json:"Role"`
	UserCharacters []interface{} `json:"UserCharacters"`
	GachaRecords   []interface{} `json:"GachaRecords"`
	UserPools      []interface{} `json:"UserPools"`
}

type CharacterRsp struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Rank      string `json:"rank"`
	IsLimited bool   `json:"is_limited"`
	IsUp      bool   `json:"is_up"`
}

type UserCharacterRsp struct {
	UserID                 uint         `json:"UserID"`
	User                   UserRsp      `json:"User"`
	CharacterID            uint         `json:"CharacterID"`
	Character              CharacterRsp `json:"Character"`
	OwnedCount             int          `json:"OwnedCount"`
	Level                  int          `json:"Level"`
	FirstAcquiredAt        time.Time    `json:"FirstAcquiredAt"`
	FirstAcquiredPool      uint         `json:"FirstAcquiredPool"`
	FirstAcquiredPullCount int          `json:"FirstAcquiredPullCount"`
}

type GachaPoolConfigRsp struct {
	ID                   uint      `json:"id"`
	PoolID               uint      `json:"pool_id"`
	SRankBaseRate        float64   `json:"s_rank_base_rate"`
	ARankBaseRate        float64   `json:"a_rank_base_rate"`
	AGuaranteeInterval   int       `json:"a_guarantee_interval"`
	SPityStart           int       `json:"s_pity_start"`
	SPityStep            float64   `json:"s_pity_step"`
	SPityEnd             int       `json:"s_pity_end"`
	LimitPity            int       `json:"limit_pity"`
	LimitRateWhenS       float64   `json:"limit_rate_when_s"`
	MaxLimitedCharacters int       `json:"max_limited_characters"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type GachaPoolCharacterRsp struct {
	ID          uint         `json:"ID"`
	PoolID      uint         `json:"PoolID"`
	CharacterID uint         `json:"CharacterID"`
	Character   CharacterRsp `json:"Character"`
	CreatedAt   time.Time    `json:"CreatedAt"`
}

type UserPoolRsp struct {
	ID          int  `json:"ID"`
	UserID      uint `json:"UserID"`
	PoolID      uint `json:"PoolID"`
	PullCount   int  `json:"PullCount"`
	LastACount  int  `json:"LastACount"`
	LastSCount  int  `json:"LastSCount"`
	LastSUp     bool `json:"LastSUp"`
	LastUpCount int  `json:"LastUpCount"`
}

type GachaPoolRsp struct {
	ID                  uint                    `json:"id"`
	Name                string                  `json:"name"`
	Description         string                  `json:"description"`
	Config              GachaPoolConfigRsp      `json:"Config"`
	GachaPoolCharacters []GachaPoolCharacterRsp `json:"GachaPoolCharacters"`
	Characters          []CharacterRsp          `json:"Characters"`
	StartAt             *time.Time              `json:"start_at"`
	EndAt               *time.Time              `json:"end_at"`
	UserPools           []UserPoolRsp           `json:"UserPools"`
	IsActive            bool                    `json:"is_active"`
}

func NewCharsGetRsp(characters []models.UserCharacter, code int, message string) CharsGetRsp {
	items := make([]UserCharacterRsp, 0, len(characters))
	for _, character := range characters {
		items = append(items, mapUserCharacter(character))
	}

	return CharsGetRsp{
		Characters: items,
		Code:       code,
		Message:    message,
	}
}

func NewPoolInfoRsp(pool models.GachaPool, code int, message string) PoolInfoRsp {
	return PoolInfoRsp{
		Pool:    mapGachaPool(pool),
		Code:    code,
		Message: message,
	}
}

func mapUserCharacter(character models.UserCharacter) UserCharacterRsp {
	return UserCharacterRsp{
		UserID:                 character.UserID,
		User:                   mapUser(character.User),
		CharacterID:            character.CharacterID,
		Character:              mapCharacter(character.Character),
		OwnedCount:             character.OwnedCount,
		Level:                  character.Level,
		FirstAcquiredAt:        character.FirstAcquiredAt,
		FirstAcquiredPool:      character.FirstAcquiredPool,
		FirstAcquiredPullCount: character.FirstAcquiredPullCount,
	}
}

func mapUser(user models.User) UserRsp {
	return UserRsp{
		UID:            user.UID,
		Username:       user.Username,
		Nickname:       user.Nickname,
		Profile:        user.Profile,
		Password:       user.Password,
		Email:          user.Email,
		Role:           user.Role,
		UserCharacters: nil,
		GachaRecords:   nil,
		UserPools:      nil,
	}
}

func mapCharacter(character models.Character) CharacterRsp {
	return CharacterRsp{
		ID:        character.ID,
		Name:      character.Name,
		Rank:      character.Rank,
		IsLimited: character.IsLimited,
		IsUp:      character.IsUp,
	}
}

func mapPoolConfig(config models.GachaPoolConfig) GachaPoolConfigRsp {
	return GachaPoolConfigRsp{
		ID:                   config.ID,
		PoolID:               config.PoolID,
		SRankBaseRate:        config.SRankBaseRate,
		ARankBaseRate:        config.ARankBaseRate,
		AGuaranteeInterval:   config.AGuaranteeInterval,
		SPityStart:           config.SPityStart,
		SPityStep:            config.SPityStep,
		SPityEnd:             config.SPityEnd,
		LimitPity:            config.LimitPity,
		LimitRateWhenS:       config.LimitRateWhenS,
		MaxLimitedCharacters: config.MaxLimitedCharacters,
		CreatedAt:            config.CreatedAt,
		UpdatedAt:            config.UpdatedAt,
	}
}

func mapGachaPool(pool models.GachaPool) GachaPoolRsp {
	poolCharacters := make([]GachaPoolCharacterRsp, 0, len(pool.GachaPoolCharacters))
	for _, poolCharacter := range pool.GachaPoolCharacters {
		poolCharacters = append(poolCharacters, GachaPoolCharacterRsp{
			ID:          poolCharacter.ID,
			PoolID:      poolCharacter.PoolID,
			CharacterID: poolCharacter.CharacterID,
			Character:   mapCharacter(poolCharacter.Character),
			CreatedAt:   poolCharacter.CreatedAt,
		})
	}

	characters := make([]CharacterRsp, 0, len(pool.Characters))
	for _, character := range pool.Characters {
		characters = append(characters, mapCharacter(character))
	}

	userPools := make([]UserPoolRsp, 0, len(pool.UserPools))
	for _, userPool := range pool.UserPools {
		userPools = append(userPools, UserPoolRsp{
			ID:          userPool.ID,
			UserID:      userPool.UserID,
			PoolID:      userPool.PoolID,
			PullCount:   userPool.PullCount,
			LastACount:  userPool.LastACount,
			LastSCount:  userPool.LastSCount,
			LastSUp:     userPool.LastSUp,
			LastUpCount: userPool.LastUpCount,
		})
	}

	return GachaPoolRsp{
		ID:                  pool.ID,
		Name:                pool.Name,
		Description:         pool.Description,
		Config:              mapPoolConfig(pool.Config),
		GachaPoolCharacters: poolCharacters,
		Characters:          characters,
		StartAt:             pool.StartAt,
		EndAt:               pool.EndAt,
		UserPools:           userPools,
		IsActive:            pool.IsActive,
	}
}

type CharFirstInfoRsp struct {
	Code                   int       `json:"code,omitempty"`
	Message                string    `json:"message"`
	FirstAcquiredAt        time.Time `json:"first_acquired_at,omitempty"`
	FirstAcquiredPool      uint      `json:"first_acquired_pool,omitempty"`
	FirstAcquiredPullCount int       `json:"first_acquired_pull_count,omitempty"`
}

type PoolConfigUpdateRsp struct {
	Code    int    `json:"code,omitempty"`
	PoolID  uint   `json:"pool_id"`
	Version uint64 `json:"version"`
	Message string `json:"message"`
}

type PoolIdsRsp struct {
	PoolIds []uint `json:"pool_ids"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}
