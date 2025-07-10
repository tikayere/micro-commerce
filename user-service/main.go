package main

import (
	"fmt"
	"net"

	examplev1 "github.com/tikayere/userservice/gen/example/v1"
	"github.com/tikayere/userservice/internal/config"
	"github.com/tikayere/userservice/internal/db"
	"github.com/tikayere/userservice/internal/handler"
	"github.com/tikayere/userservice/internal/repository"
	"github.com/tikayere/userservice/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Initialize database
	dbConn, err := db.New(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	sqlDB, err := dbConn.DB()
	if err != nil {
		logger.Fatal("Failed to get underlying database connection", zap.Error(err))
	}
	defer sqlDB.Close()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Initialize repository, service, and handler
	repo := repository.NewUserRepository(dbConn)
	svc := service.NewUserService(repo, cfg.JWTSecret, logger)
	h := handler.NewUserHandler(svc, logger)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	examplev1.RegisterUserServiceServer(grpcServer, h)
	logger.Info("gRPC server listening", zap.Int("port", cfg.GRPCPort))
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to server", zap.Error(err))
	}
}
