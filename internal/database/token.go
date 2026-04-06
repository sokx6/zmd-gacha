package database

import (
	"errors"
	"fmt"
	"time"
	"zmd-gacha/internal/models"

	"gorm.io/gorm"
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

func (db *Database) ValidateRefreshToken(uid uint, token string) (bool, error) {
	var rt models.RefreshToken
	err := db.DB.Where("uid = ? AND token = ?", uid, token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("查询刷新令牌失败: %w", err)
	}
	if rt.Expired || time.Now().After(rt.ExpiredAt) {
		return false, nil
	}
	return true, nil
}
