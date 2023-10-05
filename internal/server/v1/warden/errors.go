package warden

import "errors"

var (
	ErrTeamNotFound = errors.New("email is not registered on warden")
	ErrUserNotFound = errors.New("user not authorized")
)
