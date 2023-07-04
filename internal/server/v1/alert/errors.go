package alert

import "errors"

var ErrSubscriptionNotFound = errors.New("could not find subscription")
var ErrNoShieldSlackMetadata = errors.New("could not find slack metadata")
var ErrInvalidShieldSlackMetadata = errors.New("invalid slack metadata format")
var ErrNoShieldSlackChannel = errors.New("could not find channel with given severity")
var ErrInvalidSlackChannelFormat = errors.New("invalid channel name format")
