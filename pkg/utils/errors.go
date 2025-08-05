package utils

import "errors"

var (
	ErrorQueryString       = errors.New("error create query string")
	ErrorUserAlreadyExists = errors.New("username or email already exists")
	ErrorUserNotFound      = errors.New("user not found")
	ErrorInvalidPassword   = errors.New("invalid password")
	ErrorNotFoundRows      = errors.New("error finding rows")
)
