package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/errors"
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseAirplaneStateModeEnable(t *testing.T) {
	t.Parallel()
	data := "enabled"
	expected := &parser.AirplaneModeState{
		Enabled: true,
	}
	airplaneModeStatus := parser.NewAirplaneModeState()
	err := airplaneModeStatus.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, airplaneModeStatus)
}

func TestParseAirplaneStateModeDisable(t *testing.T) {
	t.Parallel()
	data := "disabled"
	expected := &parser.AirplaneModeState{
		Enabled: false,
	}
	airplaneModeStatus := parser.NewAirplaneModeState()
	err := airplaneModeStatus.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, airplaneModeStatus)
}

func TestParseAirplaneStateModeEnableOld(t *testing.T) {
	t.Parallel()
	data := "1"
	expected := &parser.AirplaneModeState{
		Enabled: true,
	}
	airplaneModeStatus := parser.NewAirplaneModeState()
	err := airplaneModeStatus.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, airplaneModeStatus)
}

func TestParseAirplaneStateModeDisableOld(t *testing.T) {
	t.Parallel()
	data := "0"
	expected := &parser.AirplaneModeState{
		Enabled: false,
	}
	airplaneModeStatus := parser.NewAirplaneModeState()
	err := airplaneModeStatus.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, airplaneModeStatus)
}

func TestParseAirplaneStateModeErrorAndroidVersion(t *testing.T) {
	t.Parallel()
	data := "No shell command implementation. "
	airplaneModeStatus := parser.NewAirplaneModeState()
	err := airplaneModeStatus.Parse(data)
	assert.ErrorIs(t, err, errors.ErrorMinimumAndroidVersionNotSupport)
}

func BenchmarkParseAirplaneStateModeEnable(b *testing.B) {
	data := "enabled"
	for i := 0; i < b.N; i++ {
		airplaneModeStatus := parser.NewAirplaneModeState()
		_ = airplaneModeStatus.Parse(data)
	}
}

func BenchmarkParseAirplaneStateModeDisable(b *testing.B) {
	data := "disabled"
	for i := 0; i < b.N; i++ {
		airplaneModeStatus := parser.NewAirplaneModeState()
		_ = airplaneModeStatus.Parse(data)
	}
}

func BenchmarkParseAirplaneStateModeEnableOld(b *testing.B) {
	data := "1"
	for i := 0; i < b.N; i++ {
		airplaneModeStatus := parser.NewAirplaneModeState()
		_ = airplaneModeStatus.Parse(data)
	}
}
