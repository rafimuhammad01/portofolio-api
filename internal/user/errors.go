package user

import "github.com/pkg/errors"

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrInternalServer            = errors.New("internal server error")
	ErrUsernameAlreadyExist      = errors.New("username is already taken")
	ErrInvalidUsernameOrPassword = errors.New("wrong username/password")
)
