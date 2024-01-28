package parser

import "strings"

func ParseBatteryLevel(rawBatteryLevel string) string {
	return strings.TrimSpace(rawBatteryLevel)
}
