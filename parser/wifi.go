package parser

import (
	"regexp"
	"strings"
)

type Wifi struct {
	SSID      string `json:"SSID"`
	Connected bool   `json:"connected"`
}

func ParseWifi(rawWifiInfo string) Wifi {
	wifi := Wifi{}
	isConnected := strings.Contains(rawWifiInfo, "Supplicant state: COMPLETED")
	if isConnected {
		ssidPattern := `current SSID\(s\):\{iface=wlan0,ssid="([^"]+)"\}`
		re := regexp.MustCompile(ssidPattern)
		matches := re.FindStringSubmatch(rawWifiInfo)
		if len(matches) > 1 {
			wifi.SSID = matches[1]
		}
	}
	wifi.Connected = isConnected
	return wifi
}
