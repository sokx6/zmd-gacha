package main

import (
	"zmd-gacha/internal/game/handler"
	"zmd-gacha/internal/game/middleware"
	shared_middleware "zmd-gacha/internal/middleware"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware, gachaHandler *handler.GachaHandler) {
	e.Use(shared_middleware.Logger)
	api := e.Group("/api")
	api.Use(authMiddleware.Jwt)
	// 普通登录用户可访问
	api.GET("/pull", gachaHandler.Pull)
	api.GET("/characters", userHandler.GetUserCharacters)
	api.GET("/pool", gachaHandler.GetPoolInfo)
	api.GET("/character/first_info", userHandler.GetCharFirstInfo)

}
