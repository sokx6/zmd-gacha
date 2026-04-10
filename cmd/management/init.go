package main

import (
	"fmt"
	"zmd-gacha/internal/management/config"
	"zmd-gacha/internal/management/database"
	"zmd-gacha/internal/management/handler"
	"zmd-gacha/internal/management/middleware"
	"zmd-gacha/internal/management/service"

	shared_middleware "zmd-gacha/internal/middleware"

	"github.com/labstack/echo/v4"
)

type ManagerServer struct {
	Echo *echo.Echo
	Cfg  *config.Config
}

func NewServer(cfg_path string) *ManagerServer {
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
	e.Use(shared_middleware.Logger)
	e.HTTPErrorHandler = shared_middleware.AppHTTPErrorHandler
	RegisterRoutes(e, authHandler, userHandler, authMiddleware, gachaHandler)
	return &ManagerServer{
		Echo: e,
		Cfg:  cfg,
	}
}

func (s *ManagerServer) Start() {
	addr := s.Cfg.App.Host + ":" + fmt.Sprintf("%d", s.Cfg.App.Port)

	s.Echo.Logger.Fatal(s.Echo.Start(addr))
}
