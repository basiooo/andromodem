package model

import "github.com/basiooo/andromodem/parser"

type ToggleMobileDataResponse struct {
	Success bool   `json:"success"`
	Enabled bool   `json:"enabled"`
	Error   string `json:"errors,omitempty"`
}
type DeviceNetwork struct {
	AirplaneMode bool             `json:"airplane_mode"`
	Carriers     []parser.Carrier `json:"carriers"`
	Wifi         parser.Wifi      `json:"wifi"`
}
