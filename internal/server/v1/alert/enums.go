package alert

type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "INFO"
	AlertSeverityWarning  AlertSeverity = "WARNING"
	AlertSeverityCritical AlertSeverity = "CRITICAL"
)

type ChannelCriticality string

const (
	ChannelCriticalityInfo     ChannelCriticality = "info"
	ChannelCriticalityWarning  ChannelCriticality = "warning"
	ChannelCriticalityCritical ChannelCriticality = "critical"
)
