package database

import (
	"zmd-gacha/internal/models"
)

func (db *Database) GetUserCharacters(uid uint) ([]models.UserCharacter, error) {
	var userCharacters []models.UserCharacter
	if err := db.DB.Where("user_id = ?", uid).Preload("Character").Find(&userCharacters).Error; err != nil {
		return nil, err
	}
	return userCharacters, nil
}

func (db *Database) GetCharFirstInfo(uid uint, characterId uint) (*models.UserCharacter, error) {
	var userCharacter models.UserCharacter
	if err := db.DB.Where("user_id = ? AND character_id = ?", uid, characterId).Preload("Character").First(&userCharacter).Error; err != nil {
		return nil, err
	}
	return &userCharacter, nil
}

func (db *Database) GetUserPool(uid, poolId uint) (models.UserPool, error) {
	var userPool models.UserPool
	if err := db.DB.Where("user_id = ? AND pool_id = ?", uid, poolId).First(&userPool).Error; err != nil {
		return models.UserPool{}, err
	}
	return userPool, nil
}
