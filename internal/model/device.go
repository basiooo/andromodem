package model

import "github.com/basiooo/andromodem/pkg/adb/parser"

type Device struct {
	Serial string `json:"serial"`
	Model  string `json:"model"`
	State  string `json:"state"`
}

type DeviceInfo struct {
	parser.DeviceProp `json:"prop"`
	parser.Root       `json:"root"`
	parser.Battery    `json:"battery"`
}

type DevicesResponse struct {
	Devices []Device `json:"devices"`
}

type DeviceInfoResponse struct {
	DeviceInfo `json:"device_info"`
}
