package service

import (
	"zmd-gacha/internal/database"
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/utils"
)

type GachaService struct {
	DB *database.Database
}

func NewGachaService(db *database.Database) *GachaService {
	return &GachaService{
		DB: db,
	}
}

func (gs *GachaService) PullOnce(poolId uint, uid uint) (models.Character, error) {
	cfg, user, err := gs.DB.GetPullCfg(poolId, uid)
	if err != nil {
		return models.Character{}, err
	}

	var characters []models.Character
	characters, err = gs.DB.GetCharacters(poolId)
	if err != nil {
		return models.Character{}, err
	}
	result := utils.Pull(cfg, characters, user)

	if err := gs.DB.UpdateGacha(uid, poolId, result); err != nil {
		return models.Character{}, err
	}

	return result, nil
}
