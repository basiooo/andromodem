package model

import "github.com/basiooo/andromodem/pkg/adb/parser"

type BaseToggle struct {
	Enabled bool `json:"enabled"`
}

type ToggleAirplaneModeResponse struct {
	parser.AirplaneModeStatus `json:"airplane_mode"`
}
type ToggleMobileDataResponse BaseToggle

type NetworkInfo struct {
	AirplaneMode bool             `json:"airplane_mode"`
	Apn          *parser.Apn      `json:"apn"`
	Ip           string           `json:"ip"`
	Carriers     []parser.Carrier `json:"carriers"`
	// Wifi         parsers.Wifi     `json:"wifi"`
}

type NetworkInfoResponse struct {
	NetworkInfo `json:"network_info"`
}
