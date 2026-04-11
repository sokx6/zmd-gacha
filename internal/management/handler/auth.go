package handler

import (
	"net/http"
	"zmd-gacha/internal/management/service"
	"zmd-gacha/internal/types"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

// 注册处理函数
func (h *AuthHandler) Register(c echo.Context) error {
	var req types.UserRstReq

	// 各种不合法请求
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.UserRstRsp{
			Message: "不合法请求",
		})
	} else if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, types.UserRstRsp{
			Message: "用户名或密码不能为空",
		})
	} else if len(req.Username) > 40 || len(req.Password) > 40 {
		return c.JSON(http.StatusBadRequest, types.UserRstRsp{
			Message: "用户名或密码过长",
		})
	} else if req.Email == "" {
		return c.JSON(http.StatusBadRequest, types.UserRstRsp{
			Message: "邮箱不能为空",
		})
	}

	// 注册用户并生成UID
	var err error
	var uid uint
	var role string
	if uid, role, err = h.Service.Register(req.Username, req.Password, req.Email, req.Nickname, req.Profile); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.UserRstRsp{
		Code:    http.StatusOK,
		Message: "注册成功",
		UID:     uid,
		Role:    role,
	})
}

// 登录处理函数
func (h *AuthHandler) Login(c echo.Context) error {
	var req types.UserLoginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.UserLoginRsp{
			Message: "不合法请求",
		})
	} else if (req.Username == "" && req.UID == 0 && req.Email == "") || req.Password == "" {
		return c.JSON(http.StatusBadRequest, types.UserLoginRsp{
			Message: "用户名或uid或邮箱密码不能为空",
		})
	}

	_, uid, role, err := h.Service.Login(req.Username, req.Password, req.UID, req.Email)
	if err != nil {
		return err
	}

	refreshToken, err := h.Service.GenerateRefreshToken(uid)
	if err != nil {
		return err
	}

	var accessToken string
	switch role {
	case "admin":
		accessToken, err = h.Service.GenerateAdminAccessToken(uid)
		if err != nil {
			return err
		}
	default:
		accessToken, err = h.Service.GenerateUserAccessToken(uid)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.UserLoginRsp{
		Code:         http.StatusOK,
		Message:      "登录成功",
		Role:         role,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	})
}

// 刷新令牌处理函数
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req types.TokenRefReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.TokenRefRsp{
			Message: "不合法请求",
		})
	} else if req.UID == 0 || req.RefreshToken == "" {
		return c.JSON(http.StatusBadRequest, types.TokenRefRsp{
			Message: "UID或刷新令牌不能为空",
		})
	}

	newAccessToken, newRefreshToken, err := h.Service.RefreshToken(req.UID, req.RefreshToken)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.TokenRefRsp{
		Code:         http.StatusOK,
		Message:      "刷新成功",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}
