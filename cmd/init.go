package main

import (
	"fmt"
	"zmd-gacha/internal/api"
	"zmd-gacha/internal/config"
	"zmd-gacha/internal/database"
	"zmd-gacha/internal/handler"
	"zmd-gacha/internal/middleware"
	"zmd-gacha/internal/service"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo *echo.Echo
	Cfg  *config.Config
}

func NewServer(cfg_path string) *Server {
	cfg, err := config.LoadConfig(cfg_path)
	if err != nil {
		panic(err)
	}
	if err := database.Init(cfg.Database); err != nil {
		panic(err)
	}
	db, err := database.Get()
	if err != nil {
		panic(err)
	}
	authService := service.NewAuthService(db, cfg.Auth)
	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)
	gachaService := service.NewGachaService(db)
	gachaHandler := handler.NewGachaHandler(gachaService)

	e := echo.New()
	e.HTTPErrorHandler = middleware.AppHTTPErrorHandler
	api.RegisterRoutes(e, authHandler, userHandler, authMiddleware, gachaHandler)

	return &Server{
		Echo: e,
		Cfg:  cfg,
	}
}

func (s *Server) Start() {
	addr := s.Cfg.App.Host + ":" + fmt.Sprintf("%d", s.Cfg.App.Port)

	s.Echo.Logger.Fatal(s.Echo.Start(addr))
}
