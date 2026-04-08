package types

import "zmd-gacha/internal/models"

type UserRstReq struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Profile  string `json:"profile"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserLoginReq struct {
	Username string `json:"username"`
	UID      uint   `json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenRefReq struct {
	UID          uint   `json:"uid"`
	RefreshToken string `json:"refresh_token"`
}

type ProfileUpdateReq struct {
	User models.User `json:"user"`
}

type GachaPullReq struct {
	PoolID uint `json:"pool_id"`
}

type CharCreateReq struct {
	Name      string `json:"name"`
	Rank      string `json:"rank"`
	IsLimited bool   `json:"is_limited"`
	IsUp      bool   `json:"is_up"`
}

type PoolCreateReq struct {
	Pool   models.GachaPool       `json:"pool"`
	Config models.GachaPoolConfig `json:"config"`
}

type InsertCharReq struct {
	PoolId      uint `json:"pool_id"`
	CharacterId uint `json:"character_id"`
}
