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

func (db *Database) VerifyUser(user types.UserLoginReq) (bool, uint, error) {
	var dbUser models.User
	if user.Username != "" {
		if err := db.DB.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return false, 0, types.UserNotFoundError
			}
			return false, 0, err
		}
	} else if user.Email != "" {
		if err := db.DB.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return false, 0, types.UserNotFoundError
			}
			return false, 0, err
		}
	} else if user.UID != 0 {
		if err := db.DB.Where("uid = ?", user.UID).First(&dbUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return false, 0, types.UserNotFoundError
			}
			return false, 0, err
		}
	} else {
		return false, 0, types.UserError{} //todo: 定义一个更合适的错误类型
	}

	if !utils.CheckPWD(dbUser.Password, user.Password) {
		return false, 0, types.PasswordIncorrectError
	}

	return true, dbUser.UID, nil
}
