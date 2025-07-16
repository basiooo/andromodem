package parser_test

import (
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseBattery(t *testing.T) {
	t.Parallel()
	data := `Current Battery Service state:
	AC powered: false
	USB powered: true
	Wireless powered: false
	Max charging current: 500000
	Max charging voltage: 5000000
	Charge counter: 2280
	status: 2
	health: 2
	present: true
	level: 51
	scale: 100
	voltage: 3879
	temperature: 336
	technology: Li-poly
  `
	expected := &parser.Battery{
		ACPowered:          false,
		USBPowered:         true,
		WirelessPowered:    false,
		MaxChargingCurrent: 500000,
		MaxChargingVoltage: 5000000,
		ChargeCounter:      2280,
		Status:             "Charging",
		Health:             "Good",
		Present:            true,
		Level:              51,
		Scale:              100,
		Temperature:        33.6,
		Technology:         "Li-poly",
	}
	battery := parser.NewBattery()
	err := battery.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, battery)
}

func TestParseBatteryWithoutPhysicalBattery(t *testing.T) {
	t.Parallel()
	data := `Current Battery Service state:
	AC powered: false
	USB powered: true
	Wireless powered: false
	Max charging current: 36500
	Max charging voltage: 4898888
	Charge counter: 1
	status: 5
	health: 7
	present: false
	level: 100
	scale: 100
	voltage: 4929
	temperature: 111
	technology: 
  `
	expected := &parser.Battery{
		ACPowered:          false,
		USBPowered:         true,
		WirelessPowered:    false,
		MaxChargingCurrent: 36500,
		MaxChargingVoltage: 4898888,
		ChargeCounter:      1,
		Status:             "Full",
		Health:             "Cold",
		Present:            false,
		Level:              100,
		Scale:              100,
		Temperature:        11.1,
		Technology:         "",
	}
	battery := parser.NewBattery()
	err := battery.Parse(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, battery)
}

func BenchmarkParseBattery(b *testing.B) {
	data := `Current Battery Service state:
	AC powered: false
	USB powered: true
	Wireless powered: false
	Max charging current: 500000
	Max charging voltage: 5000000
	Charge counter: 2280
	status: 2
	health: 2
	present: true
	level: 51
	scale: 100
	voltage: 3879
	temperature: 336
	technology: Li-poly
  `
	for i := 0; i < b.N; i++ {
		battery := parser.NewBattery()
		_ = battery.Parse(data)
	}
}

func BenchmarkParseBatteryWithoutPhysicalBattery(b *testing.B) {
	data := `Current Battery Service state:
	AC powered: false
	USB powered: true
	Wireless powered: false
	Max charging current: 36500
	Max charging voltage: 4898888
	Charge counter: 1
	status: 5
	health: 7
	present: false
	level: 100
	scale: 100
	voltage: 4929
	temperature: 111
	technology: 
  `
	for i := 0; i < b.N; i++ {
		battery := parser.NewBattery()
		_ = battery.Parse(data)
	}
}
