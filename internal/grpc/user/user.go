package user_server

import (
	user1 "github.com/D1Y0RBEKORIFJONOV/e-commece/protos/go/gen/protos/user"
	"google.golang.org/grpc"
	"sync"
)

type UserServer struct {
	user1.UnimplementedUserServiceServer
	userService   UserService
	mu            sync.Mutex
	statusUserMap map[string]*user1.User
}

func RegisterUserServiceServer(GRPCServer *grpc.Server, userService UserService) {
	user1.RegisterUserServiceServer(GRPCServer, &UserServer{
		userService:   userService,
		statusUserMap: make(map[string]*user1.User),
	})
}
