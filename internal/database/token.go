package database

import (
	"time"
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/types"
)

func (db *Database) StoreRefreshToken(uid uint, token string, expiredAt time.Time) error {
	return db.DB.Create(&models.RefreshToken{
		UID:       uid,
		Token:     token,
		ExpiredAt: expiredAt,
		Expired:   false,
	}).Error
}

func (db *Database) DeleteRefreshToken(uid uint, token string) error {
	return db.DB.Where("uid = ? AND token = ?", uid, token).Delete(&models.RefreshToken{}).Error
}

func (db *Database) ValidateRefreshToken(uid uint, token string) (bool, string, error) {
	var rt models.RefreshToken
	err := db.DB.Where("uid = ? AND token = ?", uid, token).First(&rt).Error
	if err != nil {
		return false, "", err
	}
	if rt.Expired || time.Now().After(rt.ExpiredAt) {
		return false, "", types.InvaildTokenError
	}

	var user models.User
	err = db.DB.Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return true, "", err
	}

	return true, user.Role, nil
}
