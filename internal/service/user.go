package service

import (
	"time"
	"zmd-gacha/internal/database"
	"zmd-gacha/internal/models"
)

type UserService struct {
	DB *database.Database
}

func NewUserService(db *database.Database) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) UpdateProfile(user models.User) error {
	return s.DB.UpdateProfile(user)
}

func (s *UserService) GetUserCharacters(uid uint) ([]models.UserCharacter, error) {
	return s.DB.GetUserCharacters(uid)
}

func (s *UserService) GetCharFirstInfo(uid uint, characterId uint) (time.Time, uint, int, error) {
	if record, err := s.DB.GetCharFirstInfo(uid, characterId); err != nil {
		return time.Time{}, 0, 0, err
	} else {
		return record.FirstAcquiredAt, record.FirstAcquiredPool, record.FirstAcquiredPullCount, nil
	}
}
