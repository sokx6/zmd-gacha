package main

import (
	"zmd-gacha/internal/management/handler"
	"zmd-gacha/internal/management/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware, gachaHandler *handler.GachaHandler) {
	e.GET("/ping", handler.Ping)
	api := e.Group("/api")

	// 公开接口
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)
	api.POST("/refresh", authHandler.Refresh)

	// 普通登录用户可访问
	auth := api.Group("/auth")
	auth.Use(authMiddleware.Jwt)
	auth.PUT("/user/me", userHandler.UpdateProfile)

	// 仅管理员可访问
	admin := api.Group("/admin")
	admin.Use(authMiddleware.AdminJwt)
	admin.POST("/characters", gachaHandler.CreateCharacter)
	admin.POST("/pools", gachaHandler.CreatePool)
	admin.POST("/pools/characters", gachaHandler.InsertCharacterToPool)
}
