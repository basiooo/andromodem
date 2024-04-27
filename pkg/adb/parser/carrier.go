package parser

import "strings"

type RawCarriers struct {
	RawConnectionsState string
	RawCarriersName     string
	RawSignalsStrength  string
}

type Carrier struct {
	Name            string `json:"name"`
	ConnectionState string `json:"connection_state"`
	SimSlot         int8   `json:"sim_slot"`
	SignalStrength  `json:"signal_strength"`
}

func NewCarriers(rawCarriers RawCarriers) *[]Carrier {
	result := parseCarriers(rawCarriers)
	return result
}

func parseCarriers(rawCarriers RawCarriers) *[]Carrier {
	var result []Carrier
	carriers := strings.Split(strings.TrimSpace(rawCarriers.RawCarriersName), ",")
	rawConnectionsState := strings.Split(strings.TrimSpace(rawCarriers.RawConnectionsState), "\n")
	signalsStrength := strings.Split(strings.TrimSpace(rawCarriers.RawSignalsStrength), "\n")
	for i, carrier := range carriers {
		if strings.TrimSpace(carrier) == "" {
			continue
		}
		simSlot := i + 1
		carrierData := Carrier{
			Name:    strings.TrimSpace(carrier),
			SimSlot: int8(simSlot),
		}
		if len(rawConnectionsState) == len(carriers) {
			mobileDataConnectionState := NewMobileDataConnectionState(rawConnectionsState[i])
			carrierData.ConnectionState = mobileDataConnectionState.ConnectionState.String()
		}
		if len(signalsStrength) == len(carriers) {
			signalStrength := NewSignalStrength(signalsStrength[i])
			carrierData.SignalStrength = *signalStrength
		}
		result = append(result, carrierData)
	}
	return &result
}
