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

	var role string

	switch username {
	case "admin":
		role = "admin"
	default:
		role = "user"
	}

	db := s.DB
	var uid uint
	switch role {
	case "user":
		uid = utils.GenerateUID()
		for ; db.UIDs[uid]; uid = utils.GenerateUID() {
		}
		db.UIDs[uid] = true
	case "admin":
		uid = 1000000000
	}

	if err = db.RegisterUser(username, hashed_pwd, email, role, uid); err != nil {
		return 0, fmt.Errorf("数据库存储用户数据失败: %w", err)
	}
	return uid, nil
}

// 用户登录服务，返回是否登录成功、用户UID、角色和错误
func (s *AuthService) Login(user types.UserLoginReq) (bool, uint, string, error) {
	db := s.DB
	if db == nil {
		var err error
		db, err = database.Get()
		if err != nil {
			return false, 0, "", types.DatabaseGetError
		}
	}
	isValid, uid, err := db.VerifyUser(user)
	if err != nil {
		return false, 0, "", err
	}

	dbUser, err := db.GetUserByUID(uid)
	if err != nil {
		return false, 0, "", err
	}

	return isValid, uid, dbUser.Role, nil
}
func (s *AuthService) Logout(uid uint, token string) error {
	return s.DB.DeleteRefreshToken(uid, token)
}

func (s *AuthService) RefreshToken(uid uint, refreshToken string) (string, string, error) {
	// 验证刷新令牌
	valid, role, err := s.DB.ValidateRefreshToken(uid, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("验证刷新令牌失败: %w", err)
	}
	if !valid {
		return "", "", fmt.Errorf("刷新令牌无效")
	}

	if err := s.DB.DeleteRefreshToken(uid, refreshToken); err != nil {
		return "", "", fmt.Errorf("删除刷新令牌失败: %w", err)
	}

	newRefreshToken, err := s.GenerateRefreshToken(uid)
	var newAccessToken string
	switch role {
	case "admin":
		newAccessToken, err = s.GenerateAdminAccessToken(uid)
		if err != nil {
			return "", "", fmt.Errorf("生成管理员访问令牌失败: %w", err)
		}
	case "user":
		newAccessToken, err = s.GenerateUserAccessToken(uid)
		if err != nil {
			return "", "", fmt.Errorf("生成用户访问令牌失败: %w", err)
		}
	default:
		return "", "", fmt.Errorf("未知的用户角色: %s", role)
	}
	return newAccessToken, newRefreshToken, nil

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

func (s *AuthService) GenerateUserAccessToken(uid uint) (string, error) {
	token, err := utils.GenerateUserAccessToken(uid, s.Cfg.AccessTokenExpire, s.Cfg.Secret)
	if err != nil {
		return "", fmt.Errorf("生成访问令牌失败: %w", err)
	}
	return token, nil
}

func (s *AuthService) GenerateAdminAccessToken(uid uint) (string, error) {
	token, err := utils.GenerateAdminAccessToken(uid, s.Cfg.AccessTokenExpire, s.Cfg.Secret)
	if err != nil {
		return "", fmt.Errorf("生成访问令牌失败: %w", err)
	}
	return token, nil
}
