package warden

import "errors"

var (
	ErrUserEmailNotFound = errors.New("user email not found")
	ErrTeamUUIDNotFound  = errors.New("team with uuid not found")
)
