package warden

import "errors"

var (
	ErrUserEmailNotFound = errors.New("user email not found")
	ErrTeamUUIDNotFound  = errors.New("team with uuid not found")
	ErrNon200            = errors.New("got non-200 when hitting warden API")
)
