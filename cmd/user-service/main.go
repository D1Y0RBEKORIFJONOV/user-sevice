package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"user-service/internal/app"
	"user-service/internal/pkg/config"
	"user-service/logger"
)

func main() {
	log := logger.SetupLogger("local")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("Starting User Service",
		slog.Any("config", cfg))

	port, err := strconv.Atoi(cfg.GRPCPort)
	if err != nil {
		log.Error(err.Error())
	}
	application := app.NewApp(log, port, *cfg, cfg.Token.AccessTTL)

	go application.GRPCServer.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop
	log.Info("received shutdown signal", slog.String("signal", sig.String()))
	application.GRPCServer.Stop()
	log.Info("shutting down server")
}
