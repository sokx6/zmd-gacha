package database

import (
	"zmd-gacha/internal/models"
)

// 查询对应池子里的所有角色
func (db *Database) GetCharacters(poolId uint) ([]models.Character, error) {
	var characters []models.Character
	if err := db.DB.Model(&models.GachaPool{ID: poolId}).Association("Characters").Find(&characters); err != nil {
		return nil, err
	}
	return characters, nil
}

// 创建角色
func (db *Database) CreateCharacter(name string, rank string) (*models.Character, error) {
	character := &models.Character{
		Name: name,
		Rank: rank,
	}
	if err := db.DB.Create(character).Error; err != nil {
		return nil, err
	}
	return character, nil
}

// 将角色添加到池子里
func (db *Database) InsertCharacterToPool(poolId uint, characterId uint, isLimited bool, isUp bool) error {
	var pool models.GachaPool
	if err := db.DB.Where("id = ?", poolId).First(&pool).Error; err != nil {
		return err
	}

	var character models.Character
	if err := db.DB.Where("id = ?", characterId).First(&character).Error; err != nil {
		return err
	}

	return db.DB.Model(&pool).Association("Characters").Append(&character)
}
