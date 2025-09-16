package errs

import "errors"

var (
	ErrInternalServer        = errors.New("internal server error")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrNicknameAlreadyExists = errors.New("nickname is already in use")
	ErrWeakPassword          = errors.New("weak password, password length must be greater than or equal to 8 and less than or equal to 16")
	ErrWrongToken            = errors.New("wrong token")
	ErrCaching               = errors.New("caching error")
)
