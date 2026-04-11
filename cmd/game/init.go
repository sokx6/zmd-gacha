package main

import (
	"context"
	"fmt"
	"zmd-gacha/internal/game/config"
	"zmd-gacha/internal/game/database"
	"zmd-gacha/internal/game/handler"
	"zmd-gacha/internal/game/middleware"
	"zmd-gacha/internal/game/service"
	shared_middleware "zmd-gacha/internal/middleware"

	"github.com/labstack/echo/v4"
)

type GameServer struct {
	Echo *echo.Echo
	Cfg  *config.Config
}

func NewServer(cfg_path string) *GameServer {
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
	authMiddleware := middleware.NewAuthMiddleware(authService)
	gachaService := service.NewGachaService(db)
	go service.StartConfigWatcher(context.Background(), cfg.Grpc.ManagementAddr, cfg.Grpc.ServerID, gachaService)
	gachaHandler := handler.NewGachaHandler(gachaService)
	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)

	e := echo.New()
	e.Use(shared_middleware.Logger)
	e.HTTPErrorHandler = shared_middleware.AppHTTPErrorHandler
	RegisterRoutes(e, userHandler, authMiddleware, gachaHandler)
	return &GameServer{
		Echo: e,
		Cfg:  cfg,
	}
}

func (s *GameServer) Start() {
	addr := s.Cfg.App.Host + ":" + fmt.Sprintf("%d", s.Cfg.App.Port)

	s.Echo.Logger.Fatal(s.Echo.Start(addr))
}
