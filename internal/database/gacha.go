package database

import (
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 查询当前所有卡池
func (db *Database) GetPools() ([]models.GachaPool, error) {
	var pools []models.GachaPool
	if err := db.DB.Preload("Config").Preload("GachaPoolCharacters").Preload("GachaPoolCharacters.Character").Find(&pools).Error; err != nil {
		return nil, err
	}
	return pools, nil
}

// 查询对应卡池的配置
func (db *Database) GetPoolCfg(poolId uint) (*models.GachaPoolConfig, error) {
	var cfg models.GachaPoolConfig
	if err := db.DB.Where("pool_id = ?", poolId).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

// 查询用户和要抽的卡池配置
func (db *Database) GetPullCfg(poolId, uid uint) (models.GachaPoolConfig, models.User, error) {
	var cfg models.GachaPoolConfig
	var user models.User
	tx := db.DB.Begin()
	if err := tx.Where("pool_id = ?", poolId).First(&cfg).Error; err != nil {
		tx.Rollback()
		return models.GachaPoolConfig{}, models.User{}, err
	}
	if err := tx.Where("uid = ?", uid).First(&user).Error; err != nil {
		tx.Rollback()
		return models.GachaPoolConfig{}, models.User{}, err
	}
	tx.Commit()
	return cfg, user, nil
}

func (db *Database) getPullCfg(tx *gorm.DB, poolId, uid uint) (models.GachaPoolConfig, models.User, error) {
	var cfg models.GachaPoolConfig
	var user models.User

	if err := tx.Where("pool_id = ?", poolId).First(&cfg).Error; err != nil {
		return models.GachaPoolConfig{}, models.User{}, err
	}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("uid = ?", uid).First(&user).Error; err != nil {
		return models.GachaPoolConfig{}, models.User{}, err
	}

	return cfg, user, nil
}

func (db *Database) getCharacters(tx *gorm.DB, poolId uint) ([]models.Character, error) {
	var characters []models.Character
	if err := tx.Model(&models.GachaPool{ID: poolId}).Association("Characters").Find(&characters); err != nil {
		return nil, err
	}
	return characters, nil
}

// PullAndUpdate 在同一事务中完成读取配置/用户、抽卡计算和落库，避免并发下保底状态错乱。
func (db *Database) PullAndUpdate(uid, poolId uint, pullCount int) ([]models.Character, error) {
	var results []models.Character

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		cfg, user, err := db.getPullCfg(tx, poolId, uid)
		if err != nil {
			return err
		}

		characters, err := db.getCharacters(tx, poolId)
		if err != nil {
			return err
		}

		switch pullCount {
		case 1:
			results = []models.Character{utils.Pull(cfg, characters, user)}
		case 10:
			results = utils.PullTen(cfg, characters, user)
		default:
			results = nil
			return gorm.ErrInvalidData
		}

		var characterIds []uint
		for _, char := range results {
			characterIds = append(characterIds, char.ID)
		}

		if err := db.createRecord(tx, uid, poolId, characterIds); err != nil {
			return err
		}
		if err := db.updateUserPity(tx, uid, results); err != nil {
			return err
		}
		if err := db.updateUserCharacter(tx, uid, characterIds); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// 创建抽卡记录
func (db *Database) createRecord(tx *gorm.DB, uid, poolId uint, characterIds []uint) error {
	var user models.User
	if err := tx.Where("uid = ?", uid).First(&user).Error; err != nil {
		return err
	}
	for i, characterId := range characterIds {
		record := &models.GachaRecord{
			UserID:      uid,
			PoolID:      poolId,
			CharacterID: characterId,
			PullCount:   user.PullCount + int(i) + 1,
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}
	}
	return nil
}

// 更新用户角色信息
func (db *Database) updateUserCharacter(tx *gorm.DB, uid uint, characterIds []uint) error {
	var newChars []models.UserCharacter
	for _, characterId := range characterIds {
		var uc models.UserCharacter
		if err := tx.Where("user_id = ? AND character_id = ?", uid, characterId).First(&uc).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				newChar := models.UserCharacter{
					UserID:      uid,
					CharacterID: characterId,
					OwnedCount:  1,
					Level:       0,
				}
				newChars = append(newChars, newChar)
			} else {
				return err
			}
		} else {
			uc.OwnedCount++
			uc.Level++
			if err := tx.Save(&uc).Error; err != nil {
				return err
			}
		}
	}
	if len(newChars) > 0 {
		if err := tx.Create(&newChars).Error; err != nil {
			return err
		}
	}
	return nil
}

// 更新用户抽卡统计
func (db *Database) updateUserPity(tx *gorm.DB, uid uint, characters []models.Character) error {
	var user models.User
	if err := tx.Where("uid = ?", uid).First(&user).Error; err != nil {
		return err
	}
	for _, character := range characters {
		user.PullCount++
		rank := character.Rank
		limited := character.IsLimited
		if rank == "S" {
			user.LastSCount = user.PullCount
			if limited {
				user.LastSLimited = true
				user.LastLimitedCount = user.PullCount
			} else {
				user.LastSLimited = false
			}
		}
		if rank == "A" {
			user.LastACount = user.PullCount
		}
	}
	if err := tx.Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

// 更新一次抽卡造成的所有数据变动
func (db *Database) UpdateGacha(uid, poolId uint, characters []models.Character) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var characterIds []uint
		for _, char := range characters {
			characterIds = append(characterIds, char.ID)
		}

		if err := db.createRecord(tx, uid, poolId, characterIds); err != nil {
			return err
		}
		if err := db.updateUserPity(tx, uid, characters); err != nil {
			return err
		}
		if err := db.updateUserCharacter(tx, uid, characterIds); err != nil {
			return err
		}

		return nil
	})
}
