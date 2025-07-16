package parser

import (
	"strconv"
)

type BatteryStatus string
type BatteryHealth string

const (
	BatteryStatusUnknown     BatteryStatus = "Unknown"
	BatteryStatusCharging    BatteryStatus = "Charging"
	BatteryStatusDisCharging BatteryStatus = "Discharging"
	BatteryStatusNotCharging BatteryStatus = "Not Charging"
	BatteryStatusFull        BatteryStatus = "Full"
)

var numericToBatteryStatus = map[string]BatteryStatus{
	"1": BatteryStatusUnknown,
	"2": BatteryStatusCharging,
	"3": BatteryStatusDisCharging,
	"4": BatteryStatusNotCharging,
	"5": BatteryStatusFull,
}

const (
	BatteryHealthUnknown            BatteryHealth = "Unknown"
	BatteryHealthGood               BatteryHealth = "Good"
	BatteryHealthOverheat           BatteryHealth = "Overheat"
	BatteryHealthDead               BatteryHealth = "Dead"
	BatteryHealthOverVoltage        BatteryHealth = "Over Voltage"
	BatteryHealthUnspecifiedFailure BatteryHealth = "Unspecified Failure"
	BatteryHealthCold               BatteryHealth = "Cold"
)

var numericToBatteryHealth = map[string]BatteryHealth{
	"1": BatteryHealthUnknown,
	"2": BatteryHealthGood,
	"3": BatteryHealthOverheat,
	"4": BatteryHealthDead,
	"5": BatteryHealthOverVoltage,
	"6": BatteryHealthUnspecifiedFailure,
	"7": BatteryHealthCold,
}

type Battery struct {
	ACPowered          bool          `json:"ac_powered"`
	USBPowered         bool          `json:"usb_powered"`
	WirelessPowered    bool          `json:"wireless_powered"`
	MaxChargingCurrent int           `json:"max_charging_current"`
	MaxChargingVoltage int           `json:"max_charging_voltage"`
	ChargeCounter      int           `json:"charge_counter"`
	Status             BatteryStatus `json:"status"`
	Health             BatteryHealth `json:"health"`
	Present            bool          `json:"present"`
	Level              uint8         `json:"level"`
	Scale              uint8         `json:"scale"`
	Temperature        float64       `json:"temperature"`
	Technology         string        `json:"technology"`
}

func NewBattery() IParser {
	return &Battery{}
}

func parseInt(value string) int {
	val, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return val
}

func (b *Battery) Parse(data string) error {
	listData := parseListData(data, []string{
		"AC powered", "USB powered", "Wireless powered",
		"Max charging current", "Max charging voltage",
		"Charge counter", "status", "health", "present",
		"level", "scale", "temperature", "technology",
	}, true)

	b.ACPowered = isTruthy(getParseListDataValue(listData, "AC powered"))
	b.USBPowered = isTruthy(getParseListDataValue(listData, "USB powered"))
	b.WirelessPowered = isTruthy(getParseListDataValue(listData, "Wireless powered"))
	b.MaxChargingCurrent = parseInt(getParseListDataValue(listData, "Max charging current"))
	b.MaxChargingVoltage = parseInt(getParseListDataValue(listData, "Max charging voltage"))
	b.ChargeCounter = parseInt(getParseListDataValue(listData, "Charge counter"))
	b.Status = numericToBatteryStatus[getParseListDataValue(listData, "status")]
	b.Health = numericToBatteryHealth[getParseListDataValue(listData, "health")]
	b.Present = isTruthy(getParseListDataValue(listData, "present"))
	b.Level = uint8(parseInt(getParseListDataValue(listData, "level")))
	b.Scale = uint8(parseInt(getParseListDataValue(listData, "scale")))
	b.Temperature = float64(parseInt(getParseListDataValue(listData, "temperature"))) / 10
	b.Technology = getParseListDataValue(listData, "technology")

	return nil
}
