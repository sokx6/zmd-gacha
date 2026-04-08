package api

import (
	"zmd-gacha/internal/handler"
	"zmd-gacha/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware, gachaHandler *handler.GachaHandler) {
	auth := e.Group("/api/auth")
	auth.GET("/ping", handler.Ping)
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)

	api := e.Group("/api")
	api.Use(authMiddleware.Jwt)

	// 普通登录用户可访问
	api.PUT("/user/me", userHandler.UpdateProfile)
	api.GET("/pull", gachaHandler.Pull)

	// 仅管理员可访问
	admin := api.Group("/admin")
	admin.Use(authMiddleware.AdminJwt)
	admin.POST("/characters", gachaHandler.CreateCharacter)
	admin.POST("/pools", gachaHandler.CreatePool)
	admin.POST("/pools/characters", gachaHandler.InsertCharacterToPool)
}
