package alert

type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "INFO"
	AlertSeverityWarning  AlertSeverity = "WARNING"
	AlertSeverityCritical AlertSeverity = "CRITICAL"
)

func toAlertSeverity(val string) AlertSeverity {
	res := AlertSeverity(val)
	switch res {
	case AlertSeverityInfo, AlertSeverityWarning, AlertSeverityCritical:
		return res
	default:
		return ""
	}
}

type ChannelCriticality string

const (
	ChannelCriticalityInfo     ChannelCriticality = "info"
	ChannelCriticalityWarning  ChannelCriticality = "warning"
	ChannelCriticalityCritical ChannelCriticality = "critical"
)

func toChannelCriticality(val string) ChannelCriticality {
	res := ChannelCriticality(val)
	switch res {
	case ChannelCriticalityInfo, ChannelCriticalityWarning, ChannelCriticalityCritical:
		return res
	default:
		return ""
	}
}
