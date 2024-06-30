package suit_tests

import (
	"context"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
	"user-service/internal/entity"
	"user-service/internal/infastructura/repository/postgresql/user"
	"user-service/internal/pkg/config"
	"user-service/internal/pkg/postgres"
	"user-service/logger"
)

type UserTesting struct {
	suite.Suite
	CleanUpFunc func()
	Repository  *user.UserRepository
}

func (s *UserTesting) SetupTest() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	pgPool, err := postgres.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	log1 := logger.SetupLogger("local")
	s.Repository = user.NewProductRepository(pgPool, log1)

	s.CleanUpFunc = func() {
	}
}

func (s *UserTesting) TearDownTest() {
	if s.CleanUpFunc != nil {
		s.CleanUpFunc()
	}
}
func TestProductTestSuite(t *testing.T) {
	suite.Run(t, new(UserTesting))
}

func (s *UserTesting) TestUser() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	userCreateReq := &entity.CreateUserReq{
		FistName:   "Diyorbek",
		LastName:   "Orifjonov",
		Email:      "diyordev3@gmail.com",
		Password:   "1234",
		SecretCode: "122334",
	}
	err := s.Repository.SaveUser(ctx, userCreateReq)
	s.NoError(err)

	register1, err := s.Repository.Register(ctx, &entity.ReadUserReq{
		Email:      userCreateReq.Email,
		SecretCode: userCreateReq.SecretCode,
	})
	s.NoError(err)

	getUserRes, err := s.Repository.GetUser(ctx, &entity.FieldValueReq{
		Field: "id",
		Value: register1.ID,
	})
	if err != entity.ErrUserDeleted || err != entity.ErrUserNotRegistered {
		s.NoError(err)
	}

	ok, err := s.Repository.CheckCodeTheSending(ctx, &entity.ReadUserReq{
		Email:      userCreateReq.Email,
		SecretCode: userCreateReq.SecretCode,
	})
	s.NoError(err)
	s.True(ok)
	s.NotNil(getUserRes)

	all, err := s.Repository.GetAllUsers(ctx, &entity.GetAllUserReq{
		Limit: 10,
	})
	s.NoError(err)
	s.NotNil(all)

	updaterUser, err := s.Repository.UpdateUser(ctx, &entity.UpdateUserReq{
		FistName: "DIYPORBEKL",
		LastName: "Orifjonov",
		Password: "1234",
		Id:       getUserRes.ID,
	})
	s.NoError(err)
	s.NotNil(updaterUser)

	err = s.Repository.DeleteUser(ctx, &entity.DeleteUserReq{
		IsHardDeleted: true,
		ID:            register1.ID,
	})
	s.NoError(err)
	select {
	case <-ctx.Done():
		s.Fail("timeout exceeded")
	default:
	}
}
