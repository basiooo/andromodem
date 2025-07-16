package parser

import (
	"encoding/json"
	"strings"
)

type RawDeviceSim struct {
	RawConnectionsState string
	RawCarriersName     string
	RawSignalsStrength  string
}
type Sim struct {
	Name            string `json:"name"`
	ConnectionState string `json:"connection_state"`
	SimSlot         uint8  `json:"sim_slot"`
	SignalStrength  `json:"signal_strength"`
}

type DeviceSim struct {
	Sims []Sim `json:"sims"`
}

func NewDeviceSim() IParser {
	return &DeviceSim{}
}
func (d *DeviceSim) Parse(rawData string) error {
	var rawDeviceSim RawDeviceSim
	err := json.Unmarshal([]byte(rawData), &rawDeviceSim)
	if err != nil {
		return err
	}
	err = d.parseDeviceSim(rawDeviceSim)
	if err != nil {
		return err
	}
	return nil
}

func (d *DeviceSim) parseDeviceSim(rawData RawDeviceSim) error {
	carriers := strings.Split(strings.TrimSpace(rawData.RawCarriersName), ",")
	rawMobileDataState := strings.Split(strings.TrimSpace(rawData.RawConnectionsState), "\n")

	rawSignalsStrength := strings.Split(strings.TrimSpace(rawData.RawSignalsStrength), ",")
	if len(rawSignalsStrength) > 2 {
		rawSignalsStrength = strings.Split(strings.TrimSpace(rawData.RawSignalsStrength), "\n")
	}

	var sims []Sim

	for i, carrier := range carriers {
		carrier = strings.TrimSpace(carrier)
		if carrier == "" {
			continue
		}

		mobileDataState := parseMobileDataState(rawMobileDataState, i)
		signalStrength := parseSignalStrength(rawSignalsStrength, i)

		simData := Sim{
			Name:            carrier,
			ConnectionState: mobileDataState.String(),
			SimSlot:         uint8(i + 1),
			SignalStrength:  *signalStrength,
		}
		sims = append(sims, simData)
	}

	d.Sims = sims
	return nil
}

func parseMobileDataState(states []string, index int) *MobileDataState {
	if index >= len(states) {
		return &MobileDataState{State: DataUnknown}
	}
	mobileState := NewMobileDataState()
	err := mobileState.Parse(states[index])
	if err != nil {
		return &MobileDataState{State: DataUnknown}
	}
	return mobileState.(*MobileDataState)
}

func parseSignalStrength(signals []string, index int) *SignalStrength {
	if index >= len(signals) {
		return &SignalStrength{}
	}
	signal := NewSignalStrength()
	err := signal.Parse(signals[index])
	if err != nil {
		return &SignalStrength{}
	}
	return signal.(*SignalStrength)
}
