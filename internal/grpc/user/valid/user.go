package valid

import (
	user1 "github.com/D1Y0RBEKORIFJONOV/e-commece/protos/go/gen/protos/user"
	"user-service/internal/entity"
)

func ValidateUserCreate(req *user1.CreateUserReq) error {
	if req.FirstName == "" {
		return entity.ErrorInvalidArguments
	}
	if req.LastName == "" {
		return entity.ErrorInvalidArguments
	}
	if req.Email == "" {
		return entity.ErrorInvalidArguments
	}
	if req.Password == "" {
		return entity.ErrorInvalidArguments
	}
	return nil
}
func ValidUserRegister(req *user1.RegisterReq) error {
	if req.SecretCode == "" {
		return entity.ErrorInvalidArguments
	}
	if req.Email == "" {
		return entity.ErrorInvalidArguments
	}
	return nil
}
func ValidateUserLogin(req *user1.LoginReq) error {
	if req.Password == "" {
		return entity.ErrorInvalidArguments
	}
	if req.Email == "" {
		return entity.ErrorInvalidArguments
	}
	return nil
}
func ValidDataGetUser(req *user1.GetUserReq) error {
	if req.Value == "" {
		return entity.ErrorInvalidArguments
	}
	if req.Field == "" {
		return entity.ErrorInvalidArguments
	}
	return nil
}

func ValidReqUpdate(req *user1.UpdatedUserReq) error {
	if req.LastName == "" && req.FirstName == "" && req.Password == "" {
		return entity.ErrReqIsEmpty
	}
	if req.Id == "" {
		return entity.ErrorInvalidArguments
	}
	return nil
}
