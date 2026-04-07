package api

import (
	"zmd-gacha/internal/handler"
	"zmd-gacha/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware, gachaHandler *handler.GachaHandler) {
	// 认证路由，不需要保护
	auth := e.Group("/api/auth")
	auth.GET("/ping", handler.Ping)
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)

	// 其他路由需要认证
	api := e.Group("/api")
	api.Use(authMiddleware.Jwt)
	api.PUT("/user/me", userHandler.UpdateProfile)
	api.GET("/pull", gachaHandler.PullOnce)
}
