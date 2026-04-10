package service

import (
	"errors"
	"net/http"
	"time"
	"zmd-gacha/internal/management/database"
	"zmd-gacha/internal/models"
	"zmd-gacha/internal/types"

	"gorm.io/gorm"
)

type UserService struct {
	DB *database.Database
}

func NewUserService(db *database.Database) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) UpdateProfile(uid uint, nickname, profile string) error {
	user := models.User{
		UID:      uid,
		Nickname: nickname,
		Profile:  profile,
	}
	if err := s.DB.UpdateProfile(user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.NewAppError(http.StatusNotFound, "未找到对应的用户", err)
		} else {
			return types.NewAppError(http.StatusInternalServerError, "更新用户信息数据库错误", err)
		}
	}
	return nil
}

func (s *UserService) GetUserCharacters(uid uint) ([]models.UserCharacter, error) {
	characters, err := s.DB.GetUserCharacters(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.NewAppError(http.StatusNotFound, "未找到对应的用户角色信息", err)
		}
		return nil, types.NewAppError(http.StatusInternalServerError, "查询用户角色信息数据库错误", err)
	}
	return characters, nil
}

func (s *UserService) GetCharFirstInfo(uid uint, characterId uint) (time.Time, uint, int, error) {
	if record, err := s.DB.GetCharFirstInfo(uid, characterId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return time.Time{}, 0, 0, types.NewAppError(http.StatusNotFound, "未找到对应的用户角色", err)
		}
		return time.Time{}, 0, 0, types.NewAppError(http.StatusInternalServerError, "查询用户角色首次获取信息数据库错误", err)
	} else {
		return record.FirstAcquiredAt, record.FirstAcquiredPool, record.FirstAcquiredPullCount, nil
	}
}
