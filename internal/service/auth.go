package service

import (
	"fmt"
	"time"

	"zmd-gacha/internal/config"
	"zmd-gacha/internal/database"
	"zmd-gacha/internal/types"
	"zmd-gacha/internal/utils"
)

type AuthService struct {
	DB  *database.Database
	Cfg config.AuthConfig
}

func NewAuthService(db *database.Database, cfg config.AuthConfig) *AuthService {
	return &AuthService{DB: db, Cfg: cfg}
}

// 用户注册函数，返回UID和错误
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

// 用户登录服务，返回是否登录成功、用户UID和错误
func (s *AuthService) Login(user types.UserLoginReq) (bool, uint, error) {
	db := s.DB
	if db == nil {
		var err error
		db, err = database.Get()
		if err != nil {
			return false, 0, types.DatabaseGetError
		}
	}
	isValid, uid, err := db.VerifyUser(user)
	if err != nil {
		return false, 0, err
	}
	return isValid, uid, nil
}

func (s *AuthService) GenerateRefreshToken(uid uint) (string, error) {
	token, err := utils.GenerateRefreshToken(s.Cfg.RefreshTokenLength)
	if err != nil {
		return "", fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	now := time.Now()
	expiredAt := now.Add(time.Duration(s.Cfg.RefreshTokenExpire) * time.Second)
	if err := s.DB.StoreRefreshToken(uid, token, expiredAt); err != nil {
		return "", fmt.Errorf("存储刷新令牌失败: %w", err)
	}
	return token, nil
}

func (s *AuthService) GenerateAccessToken(uid uint) (string, error) {
	token, err := utils.GenerateAccessToken(uid, s.Cfg.Secret)
	if err != nil {
		return "", fmt.Errorf("生成访问令牌失败: %w", err)
	}
	return token, nil
}
