package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DeviceUptime struct {
	Uptime       string `json:"uptime"`
	UptimeSecond int    `json:"uptime_second"`
}

func NewDeviceUptime() IParser {
	return &DeviceUptime{}
}

func (r *DeviceUptime) Parse(data string) error {
	times := strings.Fields(data)
	if len(times) < 1 {
		r.Uptime = "0"
		r.UptimeSecond = 0
		return nil
	}

	seconds, err := parseSeconds(times[0])
	if err != nil {
		return err
	}

	r.UptimeSecond = seconds
	r.Uptime = formatDuration(seconds)
	return nil
}

func parseSeconds(input string) (int, error) {
	if idx := strings.Index(input, "."); idx != -1 {
		input = input[:idx]
	}
	return strconv.Atoi(input)
}

func formatDuration(seconds int) string {
	duration := time.Duration(seconds) * time.Second

	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	secs := seconds % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d day", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hour", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d minute", minutes))
	}
	parts = append(parts, fmt.Sprintf("%d second", secs))

	return strings.Join(parts, ", ")
}
