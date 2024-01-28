package parser

import "strings"

func ParseAirplaneModeStatus(rawAirplaneModeStatus string) bool {
	return strings.TrimSpace(rawAirplaneModeStatus) == "enabled"
}
