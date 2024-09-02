package userStorage

import "errors"

var (
	ErrUserNotFound     = errors.New("user with this email not found")
	ErrUserAlreadyExist = errors.New("user eith this email already exist")
)
