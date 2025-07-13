package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseDeviceUptime(t *testing.T) {
	t.Parallel()
	data := "22466.63 96466.21\n"
	expected := &parser.DeviceUptime{
		Uptime:       "6 hour, 14 minute, 26 second",
		UptimeSecond: 22466,
	}
	deviceUptime := parser.NewDeviceUptime()
	err := deviceUptime.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, deviceUptime)
}
func TestParseDeviceUptimeDay(t *testing.T) {
	t.Parallel()
	data := "102466.63 96466.21\n"
	expected := &parser.DeviceUptime{
		Uptime:       "1 day, 4 hour, 27 minute, 46 second",
		UptimeSecond: 102466,
	}
	deviceUptime := parser.NewDeviceUptime()
	err := deviceUptime.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, deviceUptime)
}

func TestParseDeviceEmpty(t *testing.T) {
	t.Parallel()
	data := ""
	expected := &parser.DeviceUptime{
		Uptime:       "0",
		UptimeSecond: 0,
	}
	deviceUptime := parser.NewDeviceUptime()
	err := deviceUptime.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, deviceUptime)
}

func BenchmarkParseDeviceUptime(b *testing.B) {
	data := "22466.63 96466.21\n"
	for i := 0; i < b.N; i++ {
		deviceUptime := parser.NewDeviceUptime()
		_ = deviceUptime.Parse(data)
	}
}

func BenchmarkParseDeviceUptimeDay(b *testing.B) {
	data := "102466.63 96466.21\n"
	for i := 0; i < b.N; i++ {
		deviceUptime := parser.NewDeviceUptime()
		_ = deviceUptime.Parse(data)
	}
}

func BenchmarkParseDeviceEmpty(b *testing.B) {
	data := ""
	for i := 0; i < b.N; i++ {
		deviceUptime := parser.NewDeviceUptime()
		_ = deviceUptime.Parse(data)
	}
}
