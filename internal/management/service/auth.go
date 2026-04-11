package service

import (
	"errors"
	"net/http"
	"time"

	"zmd-gacha/internal/management/config"
	"zmd-gacha/internal/management/database"
	"zmd-gacha/internal/types"
	"zmd-gacha/internal/utils"

	"gorm.io/gorm"
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
		return 0, types.NewAppError(http.StatusInternalServerError, "密码服务错误", err)
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
		for userId, ok := db.UIDs.Load(uid); ok && uid != userId; uid = utils.GenerateUID() {
		}
		db.UIDs.Store(uid, true)
	case "admin":
		uid = 1000000000
	}

	if err = db.RegisterUser(username, hashed_pwd, email, role, uid); err != nil {
		if errors.Is(gorm.ErrDuplicatedKey, err) {
			return 0, types.NewAppError(http.StatusConflict, "用户已存在", err)
		} else {
			return 0, types.NewAppError(http.StatusInternalServerError, "数据库错误", err)
		}
	}
	return uid, nil
}

// 用户登录服务，返回是否登录成功、用户UID、角色和错误
func (s *AuthService) Login(username string, password string, uid uint, email string) (bool, uint, string, error) {
	db := s.DB
	if db == nil {
		var err error
		db, err = database.Get()
		if err != nil {
			return false, 0, "", types.NewAppError(http.StatusInternalServerError, "数据库错误", err)
		}
	}
	isValid, uid, err := db.VerifyUser(username, password, uid, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, 0, "", types.NewAppError(http.StatusUnauthorized, "用户不存在", err)
		} else {
			if err.Error() == "密码错误" {
				return false, 0, "", types.NewAppError(http.StatusUnauthorized, "密码错误", err)
			}
			return false, 0, "", types.NewAppError(http.StatusInternalServerError, "数据库错误", err)
		}
	}

	dbUser, err := db.GetUserByUID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, 0, "", types.NewAppError(http.StatusUnauthorized, "用户不存在", err)
		} else {
			return false, 0, "", types.NewAppError(http.StatusInternalServerError, "数据库错误", err)
		}
	}

	return isValid, uid, dbUser.Role, nil
}

func (s *AuthService) Logout(uid uint, token string) error {
	if err := s.DB.DeleteRefreshToken(uid, token); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.NewAppError(http.StatusUnauthorized, "刷新令牌或用户不存在", err)
		} else {
			return types.NewAppError(http.StatusInternalServerError, "数据库错误", err)
		}
	}
	return nil
}

func (s *AuthService) RefreshToken(uid uint, refreshToken string) (string, string, error) {
	// 验证刷新令牌
	valid, role, err := s.DB.ValidateRefreshToken(uid, refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if valid {
				return "", "", types.NewAppError(http.StatusUnauthorized, "用户不存在", err)
			}
			return "", "", types.NewAppError(http.StatusUnauthorized, "用户或刷新令牌不存在", err)
		} else if errors.Is(err, types.InvaildTokenError) {
			return "", "", types.NewAppError(http.StatusUnauthorized, "刷新令牌无效或已过期", err)
		} else {
			return "", "", types.NewAppError(http.StatusInternalServerError, "数据库错误", err)
		}
	}

	if err := s.DB.DeleteRefreshToken(uid, refreshToken); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", types.NewAppError(http.StatusUnauthorized, "刷新令牌或用户不存在", err)
		} else {
			return "", "", types.NewAppError(http.StatusInternalServerError, "数据库错误", err)
		}
	}

	newRefreshToken, err := s.GenerateRefreshToken(uid)
	if err != nil {
		return "", "", err
	}
	var newAccessToken string
	switch role {
	case "admin":
		newAccessToken, err = s.GenerateAdminAccessToken(uid)
		if err != nil {
			return "", "", err
		}
	case "user":
		newAccessToken, err = s.GenerateUserAccessToken(uid)
		if err != nil {
			return "", "", err
		}
	default:
		return "", "", types.NewAppError(http.StatusUnauthorized, "未知的用户角色", nil)
	}
	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) GenerateRefreshToken(uid uint) (string, error) {
	token, err := utils.GenerateRefreshToken(s.Cfg.RefreshTokenLength)
	if err != nil {
		return "", types.NewAppError(http.StatusInternalServerError, "生成刷新令牌失败", err)
	}

	now := time.Now()
	expiredAt := now.Add(time.Duration(s.Cfg.RefreshTokenExpire) * time.Second)
	if err := s.DB.StoreRefreshToken(uid, token, expiredAt); err != nil {
		return "", types.NewAppError(http.StatusInternalServerError, "存储刷新令牌失败", err)
	}
	return token, nil
}

func (s *AuthService) GenerateUserAccessToken(uid uint) (string, error) {
	token, err := utils.GenerateUserAccessToken(uid, s.Cfg.AccessTokenExpire, s.Cfg.PrivateKeyPath)
	if err != nil {
		return "", types.NewAppError(http.StatusInternalServerError, "生成访问令牌失败", err)
	}
	return token, nil
}

func (s *AuthService) GenerateAdminAccessToken(uid uint) (string, error) {
	token, err := utils.GenerateAdminAccessToken(uid, s.Cfg.AccessTokenExpire, s.Cfg.PrivateKeyPath)
	if err != nil {
		return "", types.NewAppError(http.StatusInternalServerError, "生成访问令牌失败", err)
	}
	return token, nil
}
