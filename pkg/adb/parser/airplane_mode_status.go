package parser

import "strings"

type AirplaneModeStatus struct {
	Enabled bool `json:"enabled"`
}

func NewAirplaneModeStatus(rawAirplaneModeStatus string) *AirplaneModeStatus {
	airplaneModeStatus := &AirplaneModeStatus{}
	airplaneModeStatus.Enabled = isEnabled(rawAirplaneModeStatus)
	return airplaneModeStatus
}

func isEnabled(rawAirplaneModeStatus string) bool {
	return strings.TrimSpace(rawAirplaneModeStatus) == "enabled"
}
