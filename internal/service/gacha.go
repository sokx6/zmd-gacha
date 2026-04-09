package service

import (
	"errors"
	"net/http"
	"zmd-gacha/internal/database"
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/types"

	"gorm.io/gorm"
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
		return models.Character{}, types.NewAppError(http.StatusInternalServerError, "单抽数据库错误", err)
	}
	if len(results) == 0 {
		return models.Character{}, types.NewAppError(http.StatusNotFound, "未找到抽卡结果", nil)
	}

	return results[0], nil
}

func (gs *GachaService) PullTen(poolId uint, uid uint) ([]models.Character, error) {
	results, err := gs.DB.PullAndUpdate(uid, poolId, 10)
	if err != nil {
		return nil, types.NewAppError(http.StatusInternalServerError, "十连数据库错误", err)
	}
	if len(results) < 10 {
		return nil, types.NewAppError(http.StatusNotFound, "抽卡结果丢失", nil)
	}
	return results, nil
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

func (gs *GachaService) GetPoolInfo(poolId uint) (models.GachaPool, error) {
	pool, err := gs.DB.GetPoolInfo(poolId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.GachaPool{}, types.NewAppError(http.StatusNotFound, "未找到对应的卡池", err)
		}
		return models.GachaPool{}, types.NewAppError(http.StatusInternalServerError, "查询卡池信息数据库错误", err)
	}
	return pool, nil
}
