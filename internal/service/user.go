package service

import (
	"fmt"

	"zmd-gacha/internal/database"
	"zmd-gacha/internal/types"
	"zmd-gacha/internal/utils"
)

type AuthService struct {
	DB *database.Database
}

func NewAuthService(db *database.Database) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) Register(username string, password string, email string) (uint, error) {
	hashed_pwd, err := utils.HashPWD(password)
	if err != nil {
		return 0, fmt.Errorf("密码哈希失败: %w", err)
	}

	db := s.DB
	if db == nil {
		var err error
		db, err = database.Get()
		if err != nil {
			return 0, fmt.Errorf("获取数据库实例失败: %w", err)
		}
	}

	uid := utils.GenerateUID()
	for ; db.UIDs[uid]; uid = utils.GenerateUID() {
	}

	db.UIDs[uid] = true

	if err = db.RegisterUser(username, hashed_pwd, email, uid); err != nil {
		return 0, fmt.Errorf("数据库存储用户数据失败: %w", err)
	}
	return uid, nil
}

func (s *AuthService) Login(username string, password string) (bool, error) {
	db := s.DB
	if db == nil {
		var err error
		db, err = database.Get()
		if err != nil {
			return false, types.DatabaseGetError
		}
	}

	isValid, err := db.VerifyUser(username, password)
	if err != nil {
		return false, err
	}
	return isValid, nil
}
