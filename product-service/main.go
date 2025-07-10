package main

import (
	"fmt"
	"log"
	"net"

	examplev1 "github.com/tikayere/productservice/gen/example/v1"
	"github.com/tikayere/productservice/internal/config"
	"github.com/tikayere/productservice/internal/db"
	"github.com/tikayere/productservice/internal/handler"
	"github.com/tikayere/productservice/internal/repository"
	"github.com/tikayere/productservice/internal/service"
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
		log.Fatal("Failed to load config", zap.Error(err))
	}

	// Initialize database
	dbConn, err := db.New(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	SqlDB, err := dbConn.DB()
	if err != nil {
		logger.Fatal("Failed to obtained raw connection", zap.Error(err))
	}
	defer SqlDB.Close()
	SqlDB.SetMaxIdleConns(10)
	SqlDB.SetMaxOpenConns(100)

	// Initialize repository, service and handler
	repo := repository.NewProductRepository(dbConn)
	svc := service.NewProductService(repo, logger)
	h := handler.NewProductHandler(svc, logger)
	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	fmt.Printf("Starting gRPC server at :%d", cfg.GRPCPort)
	examplev1.RegisterProductServiceServer(grpcServer, h)
	logger.Info("grPC server listening", zap.Int("port", cfg.GRPCPort))
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}
