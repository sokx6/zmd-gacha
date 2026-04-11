package database

import (
	"time"
	"zmd-gacha/internal/models"

	"gorm.io/gorm"
)

func (db *Database) getPullCfg(tx *gorm.DB, poolId, uid uint) (models.GachaPoolConfig, models.UserPool, error) {
	var cfg models.GachaPoolConfig
	var user models.UserPool

	if err := tx.Where("pool_id = ?", poolId).First(&cfg).Error; err != nil {
		return models.GachaPoolConfig{}, models.UserPool{}, err
	}
	if err := tx.Where("user_id = ? AND pool_id = ?", uid, poolId).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			user = models.UserPool{
				UserID: uid,
				PoolID: poolId,
			}
			if err := tx.Create(&user).Error; err != nil {
				return models.GachaPoolConfig{}, models.UserPool{}, err
			}
		} else {
			return models.GachaPoolConfig{}, models.UserPool{}, err
		}
	}

	return cfg, user, nil
}

func (db *Database) GetCharacters(poolId uint) ([]models.Character, error) {
	var characters []models.Character
	if err := db.DB.Model(&models.GachaPool{ID: poolId}).Association("Characters").Find(&characters); err != nil {
		return nil, err
	}
	return characters, nil
}

// PullAndUpdate 在同一事务中完成读取配置/用户、抽卡计算和落库，避免并发下保底状态错乱。
func (db *Database) UpdatePullData(uid, poolId uint, results []models.Character) error {

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		var characterIds []uint
		for _, char := range results {
			characterIds = append(characterIds, char.ID)
		}

		if err := db.createRecord(tx, uid, poolId, characterIds); err != nil {
			return err
		}
		if err := db.updateUserCharacter(tx, uid, poolId, characterIds); err != nil {
			return err
		}
		if err := db.updateUserPity(tx, uid, poolId, results); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// 创建抽卡记录
func (db *Database) createRecord(tx *gorm.DB, uid, poolId uint, characterIds []uint) error {
	var userPool models.UserPool
	if err := tx.Where("user_id = ? AND pool_id = ?", uid, poolId).First(&userPool).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userPool = models.UserPool{
				UserID: uid,
				PoolID: poolId,
			}
			if err := tx.Create(&userPool).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}
	for i, characterId := range characterIds {
		record := &models.GachaRecord{
			UserID:      uid,
			PoolID:      poolId,
			CharacterID: characterId,
			PullCount:   userPool.PullCount + int(i) + 1,
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}
	}
	return nil
}

// 更新用户角色信息
func (db *Database) updateUserCharacter(tx *gorm.DB, uid, poolId uint, characterIds []uint) error {
	if len(characterIds) == 0 {
		return nil
	}
	var userPool models.UserPool
	if err := tx.Where("user_id = ? AND pool_id = ?", uid, poolId).First(&userPool).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userPool = models.UserPool{
				UserID: uid,
				PoolID: poolId,
			}
			if err := tx.Create(&userPool).Error; err != nil {
				return err
			}
		}
	}
	pullCount := userPool.PullCount
	for i, characterId := range characterIds {
		var userCharacter models.UserCharacter
		if err := tx.Where("user_id = ? AND character_id = ?", uid, characterId).First(&userCharacter).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				newUserCharacter := models.UserCharacter{
					UserID:                 uid,
					CharacterID:            characterId,
					OwnedCount:             1,
					Level:                  0,
					FirstAcquiredAt:        time.Now(),
					FirstAcquiredPool:      poolId,
					FirstAcquiredPullCount: pullCount + i + 1,
				}
				if err := tx.Create(&newUserCharacter).Error; err != nil {
					return err
				}
			}
		} else {
			userCharacter.OwnedCount++
			if userCharacter.Level < 5 {
				userCharacter.Level++
			}
			if err := tx.Save(&userCharacter).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// 更新用户抽卡统计
func (db *Database) updateUserPity(tx *gorm.DB, uid, poolId uint, characters []models.Character) error {
	var userPool models.UserPool
	if err := tx.Where("user_id = ? AND pool_id = ?", uid, poolId).First(&userPool).Error; err != nil {
		return err
	}
	for _, character := range characters {
		userPool.PullCount++
		rank := character.Rank
		up := character.IsUp
		if rank == "S" {
			userPool.LastSCount = userPool.PullCount
			if up {
				userPool.LastSUp = true
				userPool.LastUpCount = userPool.PullCount
			} else {
				userPool.LastSUp = false
			}
		}
		if rank == "A" {
			userPool.LastACount = userPool.PullCount
		}
	}
	if err := tx.Updates(&userPool).Error; err != nil {
		return err
	}
	return nil
}

func (db *Database) GetPoolInfo(poolId uint) (models.GachaPool, error) {
	var pool models.GachaPool
	if err := db.DB.Preload("Config").Preload("Characters").First(&pool, poolId).Error; err != nil {
		return models.GachaPool{}, err
	}
	return pool, nil
}

func (db *Database) GetPoolConfigs() ([]models.GachaPoolConfig, error) {
	var configs []models.GachaPoolConfig
	if err := db.DB.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}
