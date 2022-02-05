package jwt

import "github.com/pkg/errors"

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken   = errors.New("token is invalid")
	ErrExpiredToken   = errors.New("token has expired")
	ErrIntervalServer = errors.New("internal server error")
)
