package user_server

import (
	"context"
	"errors"
	user1 "github.com/D1Y0RBEKORIFJONOV/e-commece/protos/go/gen/protos/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"user-service/internal/entity"
	"user-service/internal/grpc/user/valid"
)

func (u *UserServer) CreateUser(ctx context.Context, req *user1.CreateUserReq) (*user1.Empty, error) {
	if err := valid.ValidateUserCreate(req); err != nil {
		return nil, err
	}
	ctx1, cancel1 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel1()
	err := u.userService.CreateUser(ctx1, &entity.CreateUserReq{
		Email:    req.Email,
		Password: req.Password,
		FistName: req.FirstName,
		LastName: req.LastName,
	})
	if err != nil {
		if errors.Is(err, entity.ErrorAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, err
	}
	select {
	case <-ctx1.Done():
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	default:
	}

	return &user1.Empty{}, nil
}

func (u *UserServer) RegisterUser(ctx context.Context, req *user1.RegisterReq) (*user1.User, error) {
	err := valid.ValidUserRegister(req)
	if err != nil {
		return nil, err
	}

	ctx1, cancel1 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel1()
	user, err := u.userService.Register(ctx1, &entity.ReadUserReq{
		Email:      req.Email,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		if errors.Is(err, entity.ErrorNotFound) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, err
	}
	select {
	case <-ctx1.Done():
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	default:

	}
	return &user1.User{
		Id:        user.ID,
		FirstName: user.FistName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		DeletedAt: user.DeletedAt.Format(time.RFC3339),
	}, nil
}

func (u *UserServer) Login(ctx context.Context, req *user1.LoginReq) (*user1.LoginResponse, error) {
	err := valid.ValidateUserLogin(req)
	if err != nil {
		return nil, err
	}
	ctx1, cancel1 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel1()
	token, err := u.userService.Login(ctx1, &entity.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, entity.ErrorNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, err
	}
	select {
	case <-ctx1.Done():
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	default:

	}

	return &user1.LoginResponse{
		Token: token,
	}, nil
}

func (u *UserServer) GetUser(ctx context.Context, req *user1.GetUserReq) (*user1.User, error) {
	err := valid.ValidDataGetUser(req)
	if err != nil {
		return nil, err
	}
	ctx1, cancel1 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel1()
	user, err := u.userService.GetUser(ctx1, &entity.FieldValueReq{
		Field: req.Field,
		Value: req.Value,
	})
	if err != nil {
		if errors.Is(err, entity.ErrorNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, err
	}
	select {
	case <-ctx1.Done():
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	default:

	}
	return &user1.User{
		Id:        user.ID,
		FirstName: user.FistName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		DeletedAt: user.DeletedAt.Format(time.RFC3339),
	}, nil
}

func (u *UserServer) GetAllUser(ctx context.Context, req *user1.GetAllUserReq) (*user1.GetAllUserResponse, error) {
	ctx1, cancel1 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel1()

	users, err := u.userService.GetAllUsers(ctx1, &entity.GetAllUserReq{
		Value:  req.Value,
		Field:  req.Field,
		Limit:  req.Limit,
		Offset: req.Page,
	})

	if err != nil {
		if errors.Is(err, entity.ErrInternal) {
			return nil, status.Errorf(codes.Internal, "user service error")
		}
		return nil, err
	}

	select {
	case <-ctx1.Done():
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	default:
	}
	var allUsers []*user1.User
	for i := 0; i < len(users); i++ {
		var tempUser user1.User
		tempUser.Id = users[i].ID
		tempUser.FirstName = users[i].FistName
		tempUser.LastName = users[i].LastName
		tempUser.Email = users[i].Email
		tempUser.Password = users[i].Password
		tempUser.CreatedAt = users[i].CreatedAt.Format(time.RFC3339)
		tempUser.UpdatedAt = users[i].UpdatedAt.Format(time.RFC3339)
		tempUser.DeletedAt = users[i].DeletedAt.Format(time.RFC3339)
		allUsers = append(allUsers, &tempUser)
	}
	return &user1.GetAllUserResponse{
		Users: allUsers,
	}, nil
}

func (u *UserServer) UpdateUser(ctx context.Context, req *user1.UpdatedUserReq) (*user1.User, error) {
	err := valid.ValidReqUpdate(req)
	if err != nil {
		return nil, err
	}
	ctx1, cancel1 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel1()
	user, err := u.userService.UpdateUser(ctx1, &entity.UpdateUserReq{
		Id:       req.Id,
		FistName: req.FirstName,
		LastName: req.LastName,
		Password: req.Password,
	})

	if err != nil {
		if errors.Is(err, entity.ErrorNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		if errors.Is(err, entity.ErrorNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, err
	}

	select {
	case <-ctx1.Done():
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	default:

		return &user1.User{
			Id:        user.ID,
			FirstName: user.FistName,
			LastName:  user.LastName,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
			DeletedAt: user.DeletedAt.Format(time.RFC3339),
		}, nil
	}
}

func (u *UserServer) DeleteUser(ctx context.Context, req *user1.DeleteUserReq) (*user1.DeleteUserResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}
	err := u.userService.DeleteUser(ctx, &entity.DeleteUserReq{
		ID:            req.UserId,
		IsHardDeleted: req.IsHardDelete,
	})
	if err != nil {
		if errors.Is(err, entity.ErrorNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, err
	}
	return &user1.DeleteUserResponse{
		StatusDeleted: "User deleted",
	}, nil
}
