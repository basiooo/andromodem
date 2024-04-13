package parser

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type BatteryStatus string
type BatteryHealth string

const (
	BatteryStatusUnknown     BatteryStatus = "1"
	BatteryStatusCharging    BatteryStatus = "2"
	BatteryStatusDisCharging BatteryStatus = "3"
	BatteryStatusNotCharging BatteryStatus = "4"
	BatteryStatusFull        BatteryStatus = "5"
)

var batteryStatusStrings = map[BatteryStatus]string{
	BatteryStatusUnknown:     "Unknown",
	BatteryStatusCharging:    "Charging",
	BatteryStatusDisCharging: "Discharging",
	BatteryStatusNotCharging: "Not Charging",
	BatteryStatusFull:        "Full",
}

func (batteryStatus BatteryStatus) String() string {
	if val, ok := batteryStatusStrings[batteryStatus]; ok {
		return val
	}
	return "Unknown"
}

const (
	BatteryHealthUnknown            BatteryHealth = "1"
	BatteryHealthGood               BatteryHealth = "2"
	BatteryHealthOverheat           BatteryHealth = "3"
	BatteryHealthDead               BatteryHealth = "4"
	BatteryHealthOverVoltage        BatteryHealth = "5"
	BatteryHealthUnspecifiedFailure BatteryHealth = "6"
	BatteryHealthCold               BatteryHealth = "7"
)

var batteryHealthStrings = map[BatteryHealth]string{
	BatteryHealthUnknown:            "Unknown",
	BatteryHealthGood:               "Good",
	BatteryHealthOverheat:           "Overheat",
	BatteryHealthDead:               "Dead",
	BatteryHealthOverVoltage:        "Over Voltage",
	BatteryHealthUnspecifiedFailure: "Unspecified Failure",
	BatteryHealthCold:               "Cold",
}

func (batteryHealth BatteryHealth) String() string {
	if val, ok := batteryHealthStrings[batteryHealth]; ok {
		return val
	}
	return "Unknown"
}

type Battery struct {
	ACPowered          bool          `json:"ac_powered"`
	USBPowered         bool          `json:"usb_powered"`
	WirelessPowered    bool          `json:"wireless_powered"`
	MaxChargingCurrent string        `json:"max_charging_current"`
	MaxChargingVoltage string        `json:"max_charging_voltage"`
	ChargeCounter      string        `json:"charge_counter"`
	Status             BatteryStatus `json:"status"`
	Health             BatteryHealth `json:"health"`
	Present            bool          `json:"present"`
	Level              string        `json:"level"`
	Scale              string        `json:"scale"`
	Temperature        string        `json:"temperature"`
	Technology         string        `json:"technology"`
}

func NewBattery(rawBattery string) *Battery {
	battery := &Battery{}
	err := battery.MakeBattery(rawBattery, battery)
	if err == nil {
		battery.Status = BatteryStatus(battery.Status.String())
		battery.Health = BatteryHealth(battery.Health.String())
		temp, err := strconv.ParseFloat(battery.Temperature, 64)
		if err == nil {
			temp /= 10
			battery.Temperature = strconv.FormatFloat(temp, 'f', -1, 64)
		}
	}
	return battery
}

func (b *Battery) extractBattery(rawBattery string) map[string]interface{} {
	result := make(map[string]interface{})
	lines := strings.Split(rawBattery, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == "Current Battery Service state" {
				continue
			}
			boolResult, err := strconv.ParseBool(value)
			clean_key := strings.ReplaceAll(key, " ", "_")
			clean_key = strings.ToLower(clean_key)
			result[clean_key] = value
			if err == nil && value != "0" && value != "1" {
				result[clean_key] = boolResult
			}
		}
	}
	return result
}

func (b *Battery) MakeBattery(rawBattery string, battery *Battery) error {
	batteryInterface := b.extractBattery(rawBattery)
	batteryJson, err := json.Marshal(batteryInterface)
	if err != nil {
		return errors.New("Failed convert battery interface to json")
	}

	err = json.Unmarshal([]byte(batteryJson), &battery)
	if err != nil {
		return errors.New("Failed convert battery json to Battery struct")
	}
	return nil
}
