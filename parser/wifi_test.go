package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseWifi(t *testing.T) {
	data := `Supplicant state: COMPLETED
	current SSID(s):{iface=wlan0,ssid="Ini Wifi"}`
	expected := parser.Wifi{
		SSID:      "Ini Wifi",
		Connected: true,
	}
	actual := parser.ParseWifi(data)
	assert.Equal(t, expected, actual)
}

func TestParseWifiDisconnected(t *testing.T) {
	data := `Supplicant state: DISCONNECTED
	current SSID(s):{iface=wlan0,ssid="Ini Wifi"}`
	expected := parser.Wifi{
		Connected: false,
	}
	actual := parser.ParseWifi(data)
	assert.Equal(t, expected, actual)
}
