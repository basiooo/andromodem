package model

import "github.com/basiooo/andromodem/pkg/adb/parser"

type BaseToggleResponse struct {
	Success bool   `json:"success"`
	Enabled bool   `json:"enabled"`
	Error   string `json:"errors,omitempty"`
}
type ToggleAirplaneModeResponse BaseToggleResponse
type ToggleMobileDataResponse BaseToggleResponse

type NetworkInfo struct {
	AirplaneMode bool             `json:"airplane_mode"`
	Apn          parser.Apn       `json:"APN"`
	Ip           string           `json:"Ip"`
	Carriers     []parser.Carrier `json:"carriers"`
	// Wifi         parsers.Wifi     `json:"wifi"`
}
