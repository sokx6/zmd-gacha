package database

import (
	"time"
	"zmd-gacha/internal/models"
)

func (db *Database) StoreRefreshToken(uid uint, token string, expiredAt time.Time) error {
	return db.DB.Save(&models.RefreshToken{
		UID:       uid,
		Token:     token,
		ExpiredAt: expiredAt,
		Expired:   false,
	}).Error
}
