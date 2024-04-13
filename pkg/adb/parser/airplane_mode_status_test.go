package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseAirplaneMode(t *testing.T) {
	data := "enabled"
	expected := parser.AirplaneModeStatus{
		Enabled: true,
	}
	AirplaneModeStatus := parser.NewAirplaneModeStatus(data)
	actual := *AirplaneModeStatus
	assert.Equal(t, expected, actual)
}

func TestParseAirplaneModeDisabled(t *testing.T) {
	data := "disabled"
	expected := parser.AirplaneModeStatus{
		Enabled: false,
	}
	AirplaneModeStatus := parser.NewAirplaneModeStatus(data)
	actual := *AirplaneModeStatus
	assert.Equal(t, expected, actual)
}
