package service

import (
	"errors"
	"net/http"
	"zmd-gacha/internal/game/database"
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/types"
	"zmd-gacha/internal/utils"

	"gorm.io/gorm"
)

type GachaService struct {
	DB              *database.Database
	GachaPoolConfig []models.GachaPoolConfig
}

func NewGachaService(db *database.Database) *GachaService {
	gachaPoolConfigs, err := db.GetPoolConfigs()
	if err != nil {
		panic("无法加载卡池配置: " + err.Error())
	}
	return &GachaService{
		DB:              db,
		GachaPoolConfig: gachaPoolConfigs,
	}
}

func (gs *GachaService) getPoolConfig(poolId uint) (*models.GachaPoolConfig, error) {
	for _, cfg := range gs.GachaPoolConfig {
		if cfg.PoolID == poolId {
			return &cfg, nil
		}
	}
	return nil, errors.New("未找到对应的卡池配置")
}

func (gs *GachaService) pull(poolId uint, uid uint, pullCount int) ([]models.Character, error) {
	characters, err := gs.DB.GetCharacters(poolId)
	if err != nil {
		return nil, types.NewAppError(http.StatusInternalServerError, "获取角色数据失败", err)
	}
	userPools, err := gs.DB.GetUserPool(uid, poolId)
	if err != nil {
		return nil, types.NewAppError(http.StatusInternalServerError, "获取用户抽卡数据失败", err)
	}
	poolConfig, err := gs.getPoolConfig(poolId)
	if err != nil {
		return nil, types.NewAppError(http.StatusInternalServerError, "获取卡池配置失败", err)
	}

	switch pullCount {
	case 1:
		return []models.Character{utils.Pull(*poolConfig, characters, userPools)}, nil
	case 10:
		return utils.PullTen(*poolConfig, characters, userPools), nil
	default:
		return nil, types.NewAppError(http.StatusBadRequest, "不合法的抽卡次数", nil)
	}

}

func (gs *GachaService) PullOnce(poolId uint, uid uint) (models.Character, error) {
	results, err := gs.pull(poolId, uid, 1)
	if err != nil {
		return models.Character{}, types.NewAppError(http.StatusInternalServerError, "单抽数据库错误", err)
	}
	if len(results) == 0 {
		return models.Character{}, types.NewAppError(http.StatusNotFound, "未找到抽卡结果", nil)
	}
	if err := gs.DB.UpdatePullData(uid, poolId, results); err != nil {
		return models.Character{}, types.NewAppError(http.StatusInternalServerError, "更新抽卡数据失败", err)
	}
	return results[0], nil
}

func (gs *GachaService) PullTen(poolId uint, uid uint) ([]models.Character, error) {
	results, err := gs.pull(poolId, uid, 10)
	if err != nil {
		return nil, types.NewAppError(http.StatusInternalServerError, "十连数据库错误", err)
	}
	if len(results) < 10 {
		return nil, types.NewAppError(http.StatusNotFound, "抽卡结果丢失", nil)
	}
	if err := gs.DB.UpdatePullData(uid, poolId, results); err != nil {
		return nil, types.NewAppError(http.StatusInternalServerError, "更新抽卡数据失败", err)
	}
	return results, nil
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
