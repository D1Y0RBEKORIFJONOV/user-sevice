package user_server

import (
	"context"
	"user-service/internal/entity"
)

type UserService interface {
	CreateUser(ctx context.Context, req *entity.CreateUserReq) (err error)
	Register(ctx context.Context, req *entity.ReadUserReq) (user *entity.User, err error)
	Login(ctx context.Context, req *entity.LoginReq) (token string, err error)
	GetUser(ctx context.Context, req *entity.FieldValueReq) (user *entity.User, err error)
	GetAllUsers(ctx context.Context, req *entity.GetAllUserReq) (users []*entity.User, err error)
	UpdateUser(ctx context.Context, req *entity.UpdateUserReq) (user *entity.User, err error)
	DeleteUser(ctx context.Context, req *entity.DeleteUserReq) (err error)
}
