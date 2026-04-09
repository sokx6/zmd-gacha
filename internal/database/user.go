package database

import (
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/types"
	"zmd-gacha/internal/utils"

	"gorm.io/gorm"
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

func (db *Database) UpdateProfile(user models.User) error {
	if err := db.DB.Model(&models.User{}).Where("uid = ?", user.UID).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

// 获取对应用户的所有插卡记录
func (db *Database) GetGachaRecords(uid uint) ([]models.GachaRecord, error) {
	var records []models.GachaRecord
	if err := db.DB.Where("uid = ?", uid).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// 获取对应用户信息
func (db *Database) GetUserByUID(uid uint) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, types.UserNotFoundError
		}
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
