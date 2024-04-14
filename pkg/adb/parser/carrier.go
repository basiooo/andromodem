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
		mobileDataConnectionState := NewMobileDataConnectionState(rawConnectionsState[i])
		signalStrength := NewSignalStrength(signalsStrength[i])

		simSlot := i + 1
		carrierData := Carrier{
			Name:            strings.TrimSpace(carrier),
			ConnectionState: mobileDataConnectionState.ConnectionState.String(),
			SimSlot:         int8(simSlot),
			SignalStrength:  *signalStrength,
		}
		result = append(result, carrierData)
	}
	return &result
}
