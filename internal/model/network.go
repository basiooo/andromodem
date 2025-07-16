package model

import "github.com/basiooo/andromodem/pkg/adb_processor/parser"

type Network struct {
	AirplaneMode bool               `json:"airplane_mode"`
	IpRoutes     []parser.NetworkIp `json:"ip_routes"`
	APN          parser.Apn         `json:"apn"`
	Sims         []parser.Sim       `json:"sims"`
}
