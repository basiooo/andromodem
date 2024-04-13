package model

import "github.com/basiooo/andromodem/pkg/adb/parser"

type AirplaneModeResponse struct {
	Success bool   `json:"success"`
	Enabled bool   `json:"enabled"`
	Error   string `json:"errors,omitempty"`
}
type NetworkInfo struct {
	AirplaneMode bool             `json:"airplane_mode"`
	Apn          parser.Apn       `json:"APN"`
	Ip           string           `json:"Ip"`
	Carriers     []parser.Carrier `json:"carriers"`
	// Wifi         parsers.Wifi     `json:"wifi"`
}
