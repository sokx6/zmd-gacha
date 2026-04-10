package service

import (
	"net/http"
	"zmd-gacha/internal/management/database"
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/types"
)

type GachaService struct {
	DB *database.Database
}

func NewGachaService(db *database.Database) *GachaService {
	return &GachaService{
		DB: db,
	}
}

func (gs *GachaService) CreateCharacter(name string, rank string, isLimited bool, isUp bool) (models.Character, error) {
	character, err := gs.DB.CreateCharacter(name, rank, isLimited, isUp)
	if err != nil {
		return models.Character{}, types.NewAppError(http.StatusInternalServerError, "创建角色数据库错误", err)
	}
	return character, nil
}

func (gs *GachaService) CreatePool(pool models.GachaPool, config models.GachaPoolConfig) (uint, error) {
	poolId, err := gs.DB.CreatePool(pool, config)
	if err != nil {
		return 0, types.NewAppError(http.StatusInternalServerError, "创建卡池数据库错误", err)
	}
	return poolId, nil
}

func (gs *GachaService) InsertCharacterToPool(poolId uint, characterId uint) error {
	if err := gs.DB.InsertCharacterToPool(poolId, characterId); err != nil {
		return types.NewAppError(http.StatusInternalServerError, "插入角色到卡池错误", err)
	}
	return nil
}
