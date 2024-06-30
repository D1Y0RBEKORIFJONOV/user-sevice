package user_repository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
	"user-service/internal/entity"
	"user-service/internal/pkg/postgres"
)

type UserRepository struct {
	db        *postgres.PostgresDB
	tableName string
	log       *slog.Logger
}

func NewProductRepository(db *postgres.PostgresDB, log *slog.Logger) *UserRepository {
	return &UserRepository{
		db:        db,
		tableName: "users",
		log:       log,
	}
}

//  (
//TODO
//	UserUpdater interface {
//		UpdateUser(ctx context.Context, req *entity.UpdateUserReq) (*entity.User, error)
//	}
//)

func (repo *UserRepository) SaveUser(ctx context.Context, req *entity.CreateUserReq) error {
	const op = "userRepository.SaveUser"
	log := repo.log.With(
		slog.String("method", op))
	data := map[string]interface{}{
		"first_name":  req.FistName,
		"last_name":   req.LastName,
		"email":       req.Email,
		"password":    req.Password,
		"secret_code": req.SecretCode,
	}

	query, argc, err := repo.db.Sq.Builder.
		Insert(repo.tableName).
		SetMap(data).ToSql()
	if err != nil {
		return err
	}
	log.Info("Executing query:", query)
	_, err = repo.db.Exec(ctx, query, argc...)
	if err != nil {
		return entity.ErrorAlreadyExists
	}
	return nil
}

func (repo *UserRepository) DeleteUser(ctx context.Context, req *entity.DeleteUserReq) error {
	const op = "userRepository.DeleteUser"
	log := repo.log.With(
		slog.String("method", op))
	var (
		query string
		args  []interface{}
		err   error
	)
	if req.IsHardDeleted {
		query, args, err = repo.db.Sq.Builder.Delete(repo.tableName).
			Where(repo.db.Sq.Equal("id", req.ID)).ToSql()
	} else {
		query, args, err = repo.db.Sq.Builder.Update(repo.tableName).Set("deleted_at", time.Now()).
			Where(repo.db.Sq.Equal("id", req.ID)).ToSql()
	}
	if err != nil {
		return err
	}
	log.Info("Executing query:", query)
	_, err = repo.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (user *UserRepository) GetUser(ctx context.Context, req *entity.FieldValueReq) (*entity.User, error) {
	const op = "userRepository.GetUser"
	log := user.log.With(
		slog.String("method", op),
	)
	query, args, err := user.db.Sq.Builder.
		Select(user.SelectQuery()).
		From("users").
		Where(user.db.Sq.Equal(req.Field, req.Value)).ToSql()

	if err != nil {
		log.Error("Error building SQL query", err)
		return nil, err
	}

	log.Debug("Executing query:", slog.String("query", query))
	var user1 entity.User
	var updatedAt sql.NullTime
	var deletedAt sql.NullTime
	var isRegsster bool
	err = user.db.QueryRow(ctx, query, args...).Scan(
		&user1.ID,
		&user1.FistName,
		&user1.LastName,
		&user1.Email,
		&user1.Password,
		&user1.CreatedAt,
		&updatedAt,
		&deletedAt,
		&isRegsster,
	)
	if err != nil {

		if err.Error() == "no rows in result set" {
			log.Info("User not found", slog.String("field", req.Field), slog.String("value", req.Value))
			return nil, entity.ErrorNotFound
		}
		return nil, err
	}

	if isRegsster == false {
		return nil, entity.ErrUserNotRegistered
	}
	if deletedAt.Time.IsZero() == false {

		return nil, entity.ErrUserDeleted
	}
	user1.UpdatedAt = updatedAt.Time
	return &user1, nil
}
func (repo *UserRepository) Register(ctx context.Context, user *entity.ReadUserReq) (*entity.User, error) {
	const op = "userRepository.Register"
	log := repo.log.With(
		slog.String("method", op))
	query, args, err := repo.db.Sq.Builder.Update(repo.tableName).
		Set("is_registered", true).
		Where(repo.db.Sq.Equal("email", user.Email)).ToSql()

	if err != nil {
		return nil, err
	}
	log.Info("Executing query:", query)

	_, err = repo.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	user1, err := repo.GetUser(ctx, &entity.FieldValueReq{
		Field: "email",
		Value: user.Email,
	})
	if err != nil {
		return nil, err
	}
	return user1, nil
}

func (repo *UserRepository) CheckCodeTheSending(ctx context.Context, req *entity.ReadUserReq) (bool, error) {
	const op = "userRepository.heckCodeTheSending"
	log := repo.log.With(
		slog.String("method", op))
	query, args, err := repo.db.Sq.Builder.Select("secret_code").
		From(repo.tableName).Where(repo.db.Sq.Equal("email", req.Email)).ToSql()

	if err != nil {
		return false, err
	}
	var temp_secret_code string
	log.Info("Executing query:", slog.String("query", query))
	err = repo.db.QueryRow(ctx, query, args...).Scan(&temp_secret_code)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, entity.ErrorNotFound
		}
		return false, err
	}
	return temp_secret_code == req.SecretCode, nil
}

func (repo *UserRepository) GetAllUsers(ctx context.Context, req *entity.GetAllUserReq) ([]*entity.User, error) {
	const op = "userRepository.GetAllUsers"
	log := repo.log.With(
		slog.String("method", op))
	toSql := repo.db.Sq.Builder.Select(repo.SelectQuery()).From(repo.tableName)
	if req.Field != "" && req.Value != "" {
		toSql = toSql.Where(repo.db.Sq.Equal(req.Field, req.Value))
	}
	if req.Limit != 0 {
		toSql = toSql.Limit(uint64(req.Limit))
	}
	if req.Offset != 0 {
		toSql = toSql.Offset(uint64(req.Offset))
	}
	query, args, err := toSql.ToSql()
	if err != nil {
		return nil, err
	}
	log.Info("Executing query:", query)
	var users []*entity.User
	rows, err := repo.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user entity.User
		var updatedAt sql.NullTime
		var deletedAt sql.NullTime
		var isRegsster bool
		err = rows.Scan(
			&user.ID,
			&user.FistName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&updatedAt,
			&deletedAt,
			&isRegsster)
		if err != nil {
			return nil, err
		}
		if isRegsster == false {
			continue
		}
		if deletedAt.Time.IsZero() == false {
			continue
		}
		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time
		}
		user.DeletedAt = deletedAt.Time
		users = append(users, &user)
	}
	return users, nil
}

func updateQuery(req *entity.UpdateUserReq) map[string]interface{} {
	data := map[string]interface{}{}
	if req.LastName != "" {
		data["last_name"] = req.LastName
	}
	if req.FistName != "" {
		data["first_name"] = req.FistName
	}
	if req.Password != "" {
		data["password"] = req.Password
	}
	data["updated_at"] = time.Now()

	return data
}

func (repo *UserRepository) UpdateUser(ctx context.Context, req *entity.UpdateUserReq) (*entity.User, error) {
	const op = "userRepository.UpdateUser"
	log := repo.log.With(
		slog.String("method", op))
	query, args, err := repo.db.Sq.Builder.Update(repo.tableName).
		SetMap(updateQuery(req)).
		Where(repo.db.Sq.Equal("id", req.Id)).Suffix(repo.Returning(repo.SelectQuery())).ToSql()
	if err != nil {
		return nil, err
	}

	log.Info("Executing query:", query)
	var (
		user       entity.User
		deletedAt  sql.NullTime
		isRegsster bool
	)
	err = repo.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.FistName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
		&isRegsster)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, entity.ErrorNotFound
		}
		return nil, err
	}
	if deletedAt.Time.IsZero() == false {
		return nil, entity.ErrUserDeleted
	}
	if isRegsster == false {
		return nil, entity.ErrUserNotRegistered
	}
	return &user, nil
}
