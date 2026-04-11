package main

import (
	"fmt"
	"log"
	"net"
	"zmd-gacha/internal/management/config"
	"zmd-gacha/internal/management/database"
	management_grpc "zmd-gacha/internal/management/grpc"
	"zmd-gacha/internal/management/handler"
	"zmd-gacha/internal/management/middleware"
	"zmd-gacha/internal/management/service"
	pb "zmd-gacha/proto"

	shared_middleware "zmd-gacha/internal/middleware"

	"github.com/labstack/echo/v4"
	ggrpc "google.golang.org/grpc"
)

type ManagerServer struct {
	Echo       *echo.Echo
	Cfg        *config.Config
	GrpcServer *ggrpc.Server
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
	hub := management_grpc.NewConfigHub()
	gachaService := service.NewGachaService(db, hub)
	gachaHandler := handler.NewGachaHandler(gachaService)

	e := echo.New()
	e.Use(shared_middleware.Logger)
	e.HTTPErrorHandler = shared_middleware.AppHTTPErrorHandler
	RegisterRoutes(e, authHandler, userHandler, authMiddleware, gachaHandler)

	grpcServer := ggrpc.NewServer()
	pb.RegisterConfigSyncServiceServer(grpcServer, management_grpc.NewConfigSyncServer(hub))
	return &ManagerServer{
		Echo:       e,
		Cfg:        cfg,
		GrpcServer: grpcServer,
	}
}

func (s *ManagerServer) Start() {
	host := s.Cfg.Grpc.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := s.Cfg.Grpc.Port
	if port == 0 {
		port = 9090
	}
	grpcAddr := host + ":" + fmt.Sprintf("%d", port)
	go func() {
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Fatalf("gRPC监听失败: %v", err)
		}
		if err := s.GrpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC服务失败: %v", err)
		}
	}()

	addr := s.Cfg.App.Host + ":" + fmt.Sprintf("%d", s.Cfg.App.Port)
	s.Echo.Logger.Fatal(s.Echo.Start(addr))
}
