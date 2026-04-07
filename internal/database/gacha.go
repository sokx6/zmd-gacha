package database

import (
	"zmd-gacha/internal/models"

	"gorm.io/gorm"
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

// 创建抽卡记录
func (db *Database) CreateRecord(uid, poolId, characterId uint) error {
	var user models.User
	if err := db.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		return err
	}

	record := &models.GachaRecord{
		UserID:      uid,
		PoolID:      poolId,
		CharacterID: characterId,
		PullCount:   user.PullCount + 1,
	}

	return db.DB.Create(record).Error
}

// 更新用户角色信息
func (db *Database) UpdateUserCharacter(uid, characterId uint) error {
	var userChar models.UserCharacter
	if err := db.DB.Where("user_id = ? AND character_id = ?", uid, characterId).First(&userChar).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userChar = models.UserCharacter{
				UserID:      uid,
				CharacterID: characterId,
				OwnedCount:  1,
				Level:       0,
			}
			return db.DB.Create(&userChar).Error
		}
		return err
	}
	userChar.OwnedCount++
	if userChar.Level < 6 {
		userChar.Level++
	}
	return db.DB.Save(&userChar).Error
}

// 更新用户抽卡统计
func (db *Database) UpdateUserPity(uid uint, rank string, limited bool) error {
	var user models.User
	if err := db.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		return err
	}
	user.PullCount++
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
	return db.DB.Save(&user).Error
}

// 更新一次抽卡造成的所有数据变动
func (db *Database) UpdateGacha(uid, poolId uint, character models.Character) error {
	tx := db.DB.Begin()
	if err := db.CreateRecord(uid, poolId, character.ID); err != nil {
		tx.Rollback()
		return err
	}
	if err := db.UpdateUserPity(uid, character.Rank, character.IsLimited); err != nil {
		tx.Rollback()
		return err
	}
	if err := db.UpdateUserCharacter(uid, character.ID); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
