package main

import (
	"fmt"
	"zmd-gacha/internal/api"
	"zmd-gacha/internal/config"
	"zmd-gacha/internal/database"

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
	e := echo.New()
	api.RegisterRoutes(e)

	return &Server{
		Echo: e,
		Cfg:  cfg,
	}
}

func (s *Server) Start() {
	addr := s.Cfg.App.Host + ":" + fmt.Sprintf("%d", s.Cfg.App.Port)

	s.Echo.Logger.Fatal(s.Echo.Start(addr))
}
