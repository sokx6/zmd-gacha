package service

import (
	"net/http"
	"sync/atomic"
	"zmd-gacha/internal/management/database"
	management_grpc "zmd-gacha/internal/management/grpc"
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/types"
)

type GachaService struct {
	DB      *database.Database
	Hub     *management_grpc.ConfigHub
	version atomic.Uint64
}

func NewGachaService(db *database.Database, hub *management_grpc.ConfigHub) *GachaService {
	return &GachaService{
		DB:  db,
		Hub: hub,
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

func (gs *GachaService) UpdatePoolConfig(cfg models.GachaPoolConfig) (uint64, error) {
	if err := gs.DB.UpdatePoolConfig(cfg); err != nil {
		return 0, types.NewAppError(http.StatusInternalServerError, "更新卡池配置错误", err)
	}

	newVersion := gs.version.Add(1)
	if gs.Hub != nil {
		gs.Hub.Publish(management_grpc.ConfigUpdateEvent{Version: newVersion, Config: cfg})
	}
	return newVersion, nil
}
