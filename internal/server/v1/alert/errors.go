package alert

import "errors"

var (
	ErrInvalidAlertSeverity       = errors.New("invalid alert severity")
	ErrInvalidChannelCriticality  = errors.New("invalid channel criticality")
	ErrSubscriptionNotFound       = errors.New("could not find subscription")
	ErrNoShieldSirenNamespace     = errors.New("could not find siren's namespace from project")
	ErrNoSirenReceiver            = errors.New("could not find siren's receiver")
	ErrMultipleSirenReceiver      = errors.New("multiple siren's receivers found")
	ErrInvalidSirenReceiver       = errors.New("invalid siren's receiver type")
	ErrNoShieldGroup              = errors.New("could not find shield's group")
	ErrNoSirenParentSlackReceiver = errors.New("could not find siren's parent slack receiver")
)
