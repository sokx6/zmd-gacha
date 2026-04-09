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

func (gs *GachaService) CreateCharacter(name string, rank string, isLimited bool, isUp bool) (models.Character, error) {
	return gs.DB.CreateCharacter(name, rank, isLimited, isUp)
}

func (gs *GachaService) CreatePool(pool models.GachaPool, config models.GachaPoolConfig) (uint, error) {
	return gs.DB.CreatePool(pool, config)
}

func (gs *GachaService) InsertCharacterToPool(poolId uint, characterId uint) error {
	return gs.DB.InsertCharacterToPool(poolId, characterId)
}

func (gs *GachaService) GetPoolInfo(poolId uint) (models.GachaPool, error) {
	return gs.DB.GetPoolInfo(poolId)
}
