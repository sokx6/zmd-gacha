package middleware

import (
	"net/http"
	"zmd-gacha/internal/service"
	"zmd-gacha/internal/utils"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	service *service.AuthService
}

func NewAuthMiddleware(service *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{service: service}
}

func (am *AuthMiddleware) Jwt(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		header := c.Request().Header.Get("Authorization")
		if header == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "缺少访问令牌")
		}

		if len(header) < 7 || header[:7] != "Bearer " {
			return echo.NewHTTPError(http.StatusUnauthorized, "不合法的访问令牌")
		}

		token := header[7:]

		uid, role, err := utils.ValidateAccessToken(am.service.Cfg.Secret, token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "不合法的访问令牌")
		}
		c.Set("uid", uid)
		c.Set("role", role)
		return next(c)
	}
}

func (am *AuthMiddleware) AdminJwt(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("role") != "admin" {
			return echo.NewHTTPError(http.StatusForbidden, "管理员权限不足")
		}
		return next(c)
	}
}
