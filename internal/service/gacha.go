package service

import (
	"zmd-gacha/internal/database"
	"zmd-gacha/internal/models"
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
	results, err := gs.DB.PullAndUpdate(uid, poolId, 1)
	if err != nil {
		return models.Character{}, err
	}
	if len(results) == 0 {
		return models.Character{}, err
	}

	return results[0], nil
}

func (gs *GachaService) PullTen(poolId uint, uid uint) ([]models.Character, error) {
	results, err := gs.DB.PullAndUpdate(uid, poolId, 10)
	if err != nil {
		return nil, err
	}
	return results, nil
}
