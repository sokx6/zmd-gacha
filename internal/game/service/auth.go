package service

import (
	"zmd-gacha/internal/game/config"
	"zmd-gacha/internal/game/database"
)

type AuthService struct {
	DB  *database.Database
	Cfg config.AuthConfig
}

func NewAuthService(db *database.Database, cfg config.AuthConfig) *AuthService {
	return &AuthService{DB: db, Cfg: cfg}
}
