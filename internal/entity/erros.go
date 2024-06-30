package entity

import (
	"fmt"
)

var (
	ErrorAlreadyExists    = fmt.Errorf("entity already exists")
	ErrorNotFound         = fmt.Errorf("entity not found")
	ErrorInvalidArguments = fmt.Errorf("invalid arguments")
	ErrInternal           = fmt.Errorf("internal error")
	ErrReqIsEmpty         = fmt.Errorf("req is empty")
	ErrUserNotRegistered  = fmt.Errorf("user not registered")
	ErrUserDeleted        = fmt.Errorf("user is deleted")
)
