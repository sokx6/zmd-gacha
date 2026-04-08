package types

import (
	"zmd-gacha/internal/models"
)

type UserRstRsp struct {
	UID     uint   `json:"uid"`
	Message string `json:"message"`
}

type UserLoginRsp struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type TokenRefRsp struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type ProfileUpdateRsp struct {
	Message string `json:"message"`
}

type ErrorRsp struct {
	Message string `json:"message"`
}

type PullOnceRsp struct {
	models.Character
	Message string `json:"message"`
}

type PullTenRsp struct {
	Characters []models.Character `json:"characters"`
	Message    string             `json:"message"`
}
