package parser

import (
	"strings"

	adb "github.com/abccyz/goadb"
	adbcommands "github.com/basiooo/andromodem/adb_command"
)

type RawCarrierData struct {
	RawConnectionState string
	RawCarrierNames    string
	RawSignalStrength  string
}
type Carrier struct {
	Name               string `json:"name"`
	ConnectionState    string `json:"connection_state"`
	CellSignalStrength `json:"signal_strength"`
}

func PrepareRawCarrier(d *adb.Device) RawCarrierData {
	rawSignalStrength, _ := adbcommands.GetMobileDataSignalStrength(d)
	rawConnectionState, _ := adbcommands.GetMobileDataConnectionState(d)
	rawCarrierNames, _ := adbcommands.GetCarrierNames(d)
	rawCarrierData := RawCarrierData{
		RawSignalStrength:  rawSignalStrength,
		RawConnectionState: rawConnectionState,
		RawCarrierNames:    rawCarrierNames,
	}
	return rawCarrierData
}
func ParseCarriers(rawData RawCarrierData) []Carrier {
	signalStrength := ParseSignalStrength(rawData.RawSignalStrength)
	mobileDataConnectionState := ParseMobileDataConnectionState(rawData.RawConnectionState)
	signalStrengthLength := len(signalStrength)
	carriers := strings.Split(strings.TrimSpace(rawData.RawCarrierNames), ",")
	var carriersData []Carrier
	for i, carrier := range carriers {
		carrierData := Carrier{
			Name:            strings.TrimSpace(carrier),
			ConnectionState: mobileDataConnectionState[i].String(),
		}

		if i < signalStrengthLength {
			carrierData.CellSignalStrength = signalStrength[i]
		}

		carriersData = append(carriersData, carrierData)
	}

	return carriersData
}
