package model

import (
	"time"
)

type MonitoringMethod string

const (
	MethodWS           MonitoringMethod = "ws"
	MethodICMP         MonitoringMethod = "icmp"
	MethodHTTP         MonitoringMethod = "http"
	MethodHTTPS        MonitoringMethod = "https"
	MethodPingByDevice MonitoringMethod = "ping_by_device"
)

func (m MonitoringMethod) String() string {
	switch m {
	case MethodWS:
		return "WS"
	case MethodICMP:
		return "ICMP"
	case MethodHTTP:
		return "HTTP"
	case MethodHTTPS:
		return "HTTPS"
	case MethodPingByDevice:
		return "ICMP Ping By Device"
	default:
		return "unknown"
	}
}

type MonitoringTask struct {
	Serial            string           `json:"serial" validate:"required"`
	Host              string           `json:"host" validate:"required"`
	Method            MonitoringMethod `json:"method" validate:"required,oneof=ws http https ping_by_device icmp"`
	MaxFailures       int              `json:"max_failures" validate:"required,min=1"`
	CheckingInterval  int              `json:"checking_interval" validate:"required,min=5"`
	AirplaneModeDelay int              `json:"airplane_mode_delay" validate:"required,min=0"`
	IsActive          bool             `json:"is_active"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

type MonitoringTaskRequest struct {
	Host              string           `json:"host" validate:"required"`
	Method            MonitoringMethod `json:"method" validate:"required,oneof=ws http https ping_by_device icmp"`
	MaxFailures       int              `json:"max_failures" validate:"required,min=1"`
	CheckingInterval  int              `json:"checking_interval" validate:"required,min=5"`
	AirplaneModeDelay int              `json:"airplane_mode_delay" validate:"required,min=0"`
}

type MonitoringLog struct {
	Serial  string `json:"serial"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`

	Timestamp time.Time `json:"timestamp"`
}

type MonitoringStatus struct {
	Serial       string    `json:"serial"`
	FailureCount int       `json:"failure_count"`
	LastPingTime time.Time `json:"last_ping_time"`
	IsRunning    bool      `json:"is_running"`
	LastSuccess  bool      `json:"last_success"`
}
