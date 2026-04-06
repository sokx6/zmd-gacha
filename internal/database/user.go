package database

import (
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/types"
	"zmd-gacha/internal/utils"

	"gorm.io/gorm"
)

func (db *Database) RegisterUser(username string, password string, email string, uid uint) error {
	user := models.User{
		Username: username,
		Password: password,
		Email:    email,
		UID:      uid,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		switch err {
		case gorm.ErrDuplicatedKey:
			return types.UserExistsError
		default:
			return err
		}
	}
	return nil
}

func (db *Database) VerifyUser(username string, password string) (bool, error) {
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, types.UserNotFoundError
		}
		return false, err
	}

	if !utils.CheckPWD(user.Password, password) {
		return false, types.PasswordIncorrectError
	}

	return true, nil
}
