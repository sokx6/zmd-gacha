package database

import (
	"zmd-gacha/internal/models"
)

// 创建或者更新卡池和响应设置
func (db *Database) CreatePool(pool models.GachaPool, config models.GachaPoolConfig) (uint, error) {
	gachaPool := models.GachaPool{
		Name:        pool.Name,
		Description: pool.Description,
		StartAt:     pool.StartAt,
		EndAt:       pool.EndAt,
		IsActive:    pool.IsActive,
	}
	tx := db.DB.Begin()
	if err := tx.Save(&gachaPool).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	gachaPoolConfig := models.GachaPoolConfig{
		PoolID:               gachaPool.ID,
		SRankBaseRate:        config.SRankBaseRate,
		ARankBaseRate:        config.ARankBaseRate,
		AGuaranteeInterval:   config.AGuaranteeInterval,
		SPityStart:           config.SPityStart,
		SPityStep:            config.SPityStep,
		SPityEnd:             config.SPityEnd,
		LimitPity:            config.LimitPity,
		LimitRateWhenS:       config.LimitRateWhenS,
		MaxLimitedCharacters: config.MaxLimitedCharacters,
	}

	if err := tx.Save(&gachaPoolConfig).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	return gachaPool.ID, tx.Commit().Error
}

func (db *Database) UpdatePoolConfig(cfg models.GachaPoolConfig) error {
	return db.DB.Model(&models.GachaPoolConfig{}).
		Where("pool_id = ?", cfg.PoolID).
		Updates(map[string]any{
			"s_rank_base_rate":       cfg.SRankBaseRate,
			"a_rank_base_rate":       cfg.ARankBaseRate,
			"a_guarantee_interval":   cfg.AGuaranteeInterval,
			"s_pity_start":           cfg.SPityStart,
			"s_pity_step":            cfg.SPityStep,
			"s_pity_end":             cfg.SPityEnd,
			"limit_pity":             cfg.LimitPity,
			"limit_rate_when_s":      cfg.LimitRateWhenS,
			"max_limited_characters": cfg.MaxLimitedCharacters,
		}).Error
}

func (db *Database) GetPoolConfig(poolID uint) (models.GachaPoolConfig, error) {
	var cfg models.GachaPoolConfig
	err := db.DB.Where("pool_id = ?", poolID).First(&cfg).Error
	return cfg, err
}
