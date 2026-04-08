package handler

import (
	"fmt"
	"net/http"
	"zmd-gacha/internal/service"
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
	if uid, err = h.Service.Register(req.Username, req.Password, req.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, types.UserRstRsp{
			Message: fmt.Sprintf("注册失败: %s", err.Error()),
		})
	}

	return c.JSON(http.StatusOK, types.UserRstRsp{
		Message: "注册成功",
		UID:     uid,
	})
}

// 登录处理函数
func (h *AuthHandler) Login(c echo.Context) error {
	var req types.UserLoginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.UserLoginRsp{
			Message: "不合法请求",
		})
	} else if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, types.UserLoginRsp{
			Message: "用户名或密码不能为空",
		})
	}

	isValid, uid, err := h.Service.Login(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.UserLoginRsp{
			Message: fmt.Sprintf("登录失败: %s", err.Error()), //todo 状态码不合适
		})
	} else if !isValid {
		return c.JSON(http.StatusUnauthorized, types.UserLoginRsp{
			Message: "用户名或密码错误",
		})
	}

	refreshToken, err := h.Service.GenerateRefreshToken(uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.UserLoginRsp{
			Message: "生成刷新令牌失败",
		})
	}
	accessToken, err := h.Service.GenerateUserAccessToken(uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.UserLoginRsp{
			Message: "生成访问令牌失败",
		})
	}
	return c.JSON(http.StatusOK, types.UserLoginRsp{
		Message:      "登录成功",
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
		return c.JSON(http.StatusUnauthorized, types.TokenRefRsp{
			Message: "刷新令牌失败",
		})
	}

	return c.JSON(http.StatusOK, types.TokenRefRsp{
		Message:      "刷新成功",
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}
