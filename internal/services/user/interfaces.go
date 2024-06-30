package user_service

import (
	"context"
	"user-service/internal/entity"
)

type (
	UserSaver interface {
		SaveUser(ctx context.Context, req *entity.CreateUserReq) error
		Register(ctx context.Context, req *entity.ReadUserReq) (*entity.User, error)
	}
	UserProvider interface {
		GetUser(ctx context.Context, req *entity.FieldValueReq) (*entity.User, error)
		GetAllUsers(ctx context.Context, req *entity.GetAllUserReq) ([]*entity.User, error)
		CheckCodeTheSending(ctx context.Context, req *entity.ReadUserReq) (bool, error)
	}
	UserDeleter interface {
		DeleteUser(ctx context.Context, req *entity.DeleteUserReq) error
	}
	UserUpdater interface {
		UpdateUser(ctx context.Context, req *entity.UpdateUserReq) (*entity.User, error)
	}
)
