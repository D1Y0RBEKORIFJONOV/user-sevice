package app

import (
	"log/slog"
	"time"
	grpcapp "user-service/internal/app/grpc"
	user_repository "user-service/internal/infastructura/repository/postgresql/user"
	"user-service/internal/pkg/config"
	"user-service/internal/pkg/postgres"
	user_service "user-service/internal/services/user"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(logger *slog.Logger,
	grpcPort int, configStr config.Config,
	tokenTTL time.Duration) *App {
	db, err := postgres.New(&configStr)
	if err != nil {
		panic(err)
	}
	storage := user_repository.NewProductRepository(db, logger)
	user := user_service.NewUser(logger, storage, storage, storage, storage, tokenTTL)
	GRPCApp := grpcapp.NewApp(logger, grpcPort, user)
	return &App{
		GRPCServer: GRPCApp,
	}
}
