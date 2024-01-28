package model

import "github.com/basiooo/andromodem/parser"

type DeviceInfo struct {
	parser.DeviceProps
	Root         parser.Root `json:"root"`
	BatteryLevel string      `json:"battery_level"`
}
