package entity

import "time"

type User struct {
	ID        string    `json:"id"`
	FistName  string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type CreateUserReq struct {
	FistName   string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	SecretCode string `json:"secret_code"`
}
type ReadUserReq struct {
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}
type UpdateUserReq struct {
	Id       string `json:"id"`
	FistName string `json:"first_name"`
	LastName string `json:"last_name"`
	Password string `json:"password"`
}

type DeleteUserReq struct {
	ID            string `json:"id"`
	IsHardDeleted bool   `json:"is_hard_deleted"`
}
type FieldValueReq struct {
	Field string `json:"field"`
	Value string `json:"value"`
}
type GetAllUserReq struct {
	Field  string `json:"field_value"`
	Value  string `json:"value_value"`
	Offset int64  `json:"offset_value"`
	Limit  int64  `json:"limit_value"`
}
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
