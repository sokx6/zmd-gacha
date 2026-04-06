package api

import (
	"zmd-gacha/internal/handler"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/ping", handler.Ping)
	e.POST("/register", handler.Register)
}
