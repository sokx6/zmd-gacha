package api

import (
	"zmd-gacha/internal/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, authHandler *handler.AuthHandler) {
	e.GET("/ping", handler.Ping)
	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)
	e.POST("/refresh", authHandler.Refresh)
}
