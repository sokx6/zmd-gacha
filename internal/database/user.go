package database

import (
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/types"

	"gorm.io/gorm"
)

func (db *Database) RegisterUser(username string, password string, email string) error {
	user := models.User{
		Username: username,
		Password: password,
		Email:    email,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		switch err {
		case gorm.ErrDuplicatedKey:
			return types.UserExistsError
		default:
			return types.DatabaseDefaultError
		}
	}
	return nil
}
