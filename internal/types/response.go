package types

import (
	"time"
	"zmd-gacha/internal/models"
)

type UserRstRsp struct {
	Code    int    `json:"code,omitempty"`
	UID     uint   `json:"uid"`
	Message string `json:"message"`
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
	Characters []models.UserCharacter `json:"characters"`
	Code       int                    `json:"code,omitempty"`
	Message    string                 `json:"message"`
}

type PoolInfoRsp struct {
	Pool    models.GachaPool `json:"pool"`
	Code    int              `json:"code,omitempty"`
	Message string           `json:"message"`
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
