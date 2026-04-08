package service

import (
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
