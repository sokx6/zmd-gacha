package database

import (
	"fmt"
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/utils"
)

func (db *Database) RegisterUser(username string, password string, email string, role string, uid uint) error {
	user := models.User{
		Username: username,
		Password: password,
		Email:    email,
		Role:     role,
		UID:      uid,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (db *Database) VerifyUser(username string, password string, uid uint, email string) (bool, uint, error) {
	var dbUser models.User
	if username != "" {
		if err := db.DB.Where("username = ?", username).First(&dbUser).Error; err != nil {
			return false, 0, err
		}
	} else if email != "" {
		if err := db.DB.Where("email = ?", email).First(&dbUser).Error; err != nil {
			return false, 0, err
		}
	} else if uid != 0 {
		if err := db.DB.Where("uid = ?", uid).First(&dbUser).Error; err != nil {
			return false, 0, err
		}
	} else {
		return false, 0, fmt.Errorf("缺少登录标识符")
	}

	if !utils.CheckPWD(dbUser.Password, password) {
		return false, 0, fmt.Errorf("密码错误")
	}

	return true, dbUser.UID, nil
}

func (db *Database) UpdateProfile(user models.User) error {
	if err := db.DB.Model(&models.User{}).Where("uid = ?", user.UID).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

// 获取对应用户信息
func (db *Database) GetUserByUID(uid uint) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

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
