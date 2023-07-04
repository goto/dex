package alert

import "errors"

var (
	ErrSubscriptionNotFound       = errors.New("could not find subscription")
	ErrNoShieldSlackMetadata      = errors.New("could not find slack metadata")
	ErrInvalidShieldSlackMetadata = errors.New("invalid slack metadata format")
	ErrNoShieldSlackChannel       = errors.New("could not find channel with given severity")
	ErrInvalidSlackChannelFormat  = errors.New("invalid channel name format")
)
