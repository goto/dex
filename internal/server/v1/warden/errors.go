package warden

import "errors"

var (
	ErrEmailNotOnWarden = errors.New("email is not registered on warden")
)
