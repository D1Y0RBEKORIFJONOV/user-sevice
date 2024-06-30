package user_service

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strconv"
	"time"
	"user-service/internal/entity"
	"user-service/internal/pkg/email"
	"user-service/internal/pkg/tokens"
)

func (u *User) CreateUser(ctx context.Context, req *entity.CreateUserReq) (err error) {
	const op = "user_service.CreateUser"
	log := u.log.With(
		slog.String("method-addr", op),
		slog.String("Req-email", req.Email),
		slog.String("TIME:", time.Now().UTC().Format(time.RFC3339)),
	)
	log.Info("Sending secret code to email")
	secret_code, err := email.SenSecretCode([]string{req.Email})
	if err != nil {
		return errors.Wrap(err, "failed to generate secret code")
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Info("failed to generate password hash")
		return errors.Wrap(err, "failed to generate password hash")
	}
	log.Info("CreateUser called")

	err = u.saver.SaveUser(ctx, &entity.CreateUserReq{
		Email:      req.Email,
		Password:   string(passHash),
		FistName:   req.FistName,
		LastName:   req.LastName,
		SecretCode: secret_code,
	})
	if err != nil {
		return errors.Wrap(err, op)
	}
	return nil
}
func (u *User) Register(ctx context.Context, req *entity.ReadUserReq) (user *entity.User, err error) {
	const op = "user_service.ReadUser"
	log := u.log.With(
		slog.String("method-addr", op),
		slog.String("Req-email", req.Email))

	log.Info("Sending secret code to email")
	ok, err := u.provider.CheckCodeTheSending(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	if !ok {
		return nil, errors.Wrap(errors.New("incorrect secret code input"), op)
	}
	log.Info("starting registration ")
	user, err = u.saver.Register(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return user, nil
}

func (u *User) Login(ctx context.Context, req *entity.LoginReq) (token string, err error) {
	const op = "user_service.Login"
	log := u.log.With(
		slog.String("method-addr", op),
		slog.String("Req-email", req.Email))

	log.Info("Login called")
	user, err := u.provider.GetUser(ctx, &entity.FieldValueReq{
		Value: req.Email,
		Field: "email",
	})
	if err != nil {
		return "", errors.Wrap(err, op)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.Wrap(err, op)
	}

	log.Info("Generating token")
	token, err = tokens.NewToken(user, u.tokenTTL)
	if err != nil {
		return "", errors.Wrap(err, op)
	}
	return token, nil
}

func (u *User) GetUser(ctx context.Context, req *entity.FieldValueReq) (user *entity.User, err error) {
	const op = "user_service.GetUser"
	log := u.log.With(
		slog.String("method-addr", op),
		slog.String("Field", req.Field))

	usr, err := u.provider.GetUser(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	log.Info("Retrieving user")
	return usr, nil
}

func (u *User) GetAllUsers(ctx context.Context, req *entity.GetAllUserReq) (users []*entity.User, err error) {
	const op = "user_service.GetAllUsers"
	log := u.log.With(
		slog.String("method-addr", op),
		slog.String("Limit", strconv.FormatInt(req.Limit, 10)),
		slog.String("Offset", strconv.FormatInt(req.Offset, 10)),
		slog.String("Field", req.Field))

	log.Info("Retrieving users")
	users, err = u.provider.GetAllUsers(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return users, nil
}

func (u *User) UpdateUser(ctx context.Context, req *entity.UpdateUserReq) (user *entity.User, err error) {
	const op = "user_service.UpdateUser"
	log := u.log.With(
		slog.String("method-addr", op))

	log.Info("Calling Update method")

	user, err = u.updater.UpdateUser(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	return user, nil
}

func (u *User) DeleteUser(ctx context.Context, req *entity.DeleteUserReq) (err error) {
	const op = "user_service.DeleteUser"
	log := u.log.With(
		slog.String("method-addr", op),
		slog.Bool("IsHardDelete", req.IsHardDeleted))

	log.Info("Calling Delete method")
	err = u.deleter.DeleteUser(ctx, req)
	if err != nil {
		log.Info("Oeoeeooeoeo", err)
		return errors.Wrap(err, op)
	}
	return nil
}
