package alert

import "errors"

var (
	ErrInvalidAlertSeverity      = errors.New("invalid alert severity")
	ErrInvalidChannelCriticality = errors.New("invalid channel criticality")
	ErrSubscriptionNotFound      = errors.New("could not find subscription")
	ErrNoShieldSlackChannel      = errors.New("could not find channel with given severity")
	ErrNoShieldSirenNamespace    = errors.New("could not find siren's namespace from project")
)
