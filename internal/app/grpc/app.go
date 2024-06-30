package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	user_server "user-service/internal/grpc/user"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, port int, user user_server.UserService) *App {
	grpcServer := grpc.NewServer()
	user_server.RegisterUserServiceServer(grpcServer, user)
	reflection.Register(grpcServer)
	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}
}
func (a *App) Run() error {
	const op = "grpcapp.App.Run"
	log := a.log.With(
		slog.String("method", op),
		slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	log.Info("starting gRPC server on port", a.port)
	err = a.gRPCServer.Serve(l)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}
func (a *App) Stop() {
	const op = "grpcapp.App.Stop"
	log := a.log.With(
		slog.String("method", op),
		slog.Int("port", a.port))
	log.Info("stopping gRPC server on port", a.port)
	a.gRPCServer.GracefulStop()
}
