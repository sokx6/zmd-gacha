package database

import (
	"errors"
	"net/http"
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/types"

	"gorm.io/gorm"
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

// 查询对应池子里的所有角色
func (db *Database) GetCharacters(poolId uint) ([]models.Character, error) {
	var characters []models.Character
	if err := db.DB.Model(&models.GachaPool{ID: poolId}).Association("Characters").Find(&characters); err != nil {
		return nil, err
	}
	return characters, nil
}

// 创建角色
func (db *Database) CreateCharacter(name string, rank string, isLimited bool, isUp bool) (models.Character, error) {
	character := models.Character{
		Name:      name,
		Rank:      rank,
		IsLimited: isLimited,
		IsUp:      isUp,
	}
	if err := db.DB.Create(&character).Error; err != nil {
		return models.Character{}, err
	}
	return character, nil
}

// 将角色添加到池子里
func (db *Database) InsertCharacterToPool(poolId uint, characterId uint) error {

	tx := db.DB.Begin()
	// 验证是否已经存在对应的池子和角色
	var gachaPoolCharacter models.GachaPoolCharacter
	if err := tx.Where("pool_id = ? AND character_id = ?", poolId, characterId).First(&gachaPoolCharacter).Error; err == nil {
		tx.Rollback()
		return nil // 已存在，无需添加
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}
	// 查询池子和对应配置
	var pool models.GachaPool
	if err := tx.Where("id = ?", poolId).Preload("Config").First(&pool).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.NewAppError(http.StatusNotFound, "未找到对应的卡池", err)
		}
		return err
	}
	// 查询角色
	var character models.Character
	if err := tx.Where("id = ?", characterId).First(&character).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return types.NewAppError(http.StatusNotFound, "未找到对应的角色", err)
		}
		tx.Rollback()
		return err
	}
	//获取所有角色
	var gachaPoolCharacters []models.GachaPoolCharacter

	// 限定角色的添加逻辑（有上限）
	if character.IsLimited {
		if err := tx.Model(&models.GachaPoolCharacter{}).
			Joins("JOIN characters ON characters.id = gacha_pool_characters.character_id").
			Where("gacha_pool_characters.pool_id = ? AND characters.is_limited = ?", poolId, true).
			Order("gacha_pool_characters.created_at ASC").
			Find(&gachaPoolCharacters).Error; err != nil {
			tx.Rollback()
			return err
		}
		if pool.Config.MaxLimitedCharacters <= 0 {
			tx.Rollback()
			return nil
		} else if len(gachaPoolCharacters) >= pool.Config.MaxLimitedCharacters {
			gachaPoolCharacter := gachaPoolCharacters[0]
			gachaPoolCharacter.CharacterID = character.ID
			if err := tx.Save(&gachaPoolCharacter).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			gachaPoolCharacter := models.GachaPoolCharacter{
				PoolID:      poolId,
				CharacterID: characterId,
			}
			if err := tx.Create(&gachaPoolCharacter).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	} else {
		// 非限定角色的添加逻辑（无上限）
		gachaPoolCharacter := models.GachaPoolCharacter{
			PoolID:      poolId,
			CharacterID: characterId,
		}
		if err := tx.Create(&gachaPoolCharacter).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// 从池子里删除角色
func (db *Database) RemoveCharFromPool(poolId uint, characterId uint) error {
	var gachaPoolCharacter models.GachaPoolCharacter
	tx := db.DB.Begin()
	if err := tx.Where("pool_id = ? AND character_id = ?", poolId, characterId).First(&gachaPoolCharacter).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&gachaPoolCharacter).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error

}

// 从所有池子中删除角色
func (db *Database) DeleteCharacter(characterId uint) error {
	var character models.Character
	tx := db.DB.Begin()
	if err := tx.Where("id = ?", characterId).First(&character).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Delete(&character).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
