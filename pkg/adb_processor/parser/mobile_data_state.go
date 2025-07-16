package parser

import (
	"strconv"
)

type State int8

const (
	DataUnknown      State = -1
	DataDisconnected State = 0
	DataConnecting   State = 1
	DataConnected    State = 2
	DataSuspended    State = 3
)

var stateString = map[State]string{
	DataUnknown:      "Unknown",
	DataDisconnected: "Disconnected",
	DataConnecting:   "Connecting",
	DataConnected:    "Connected",
	DataSuspended:    "Suspended",
}

func (state State) String() string {
	if val, ok := stateString[state]; ok {
		return val
	}
	return "Unknown"
}

type MobileDataState struct {
	State
}

func NewMobileDataState() IParser {
	return &MobileDataState{}
}
func (m *MobileDataState) Parse(rawData string) error {
	state, err := strconv.Atoi(rawData)
	if err != nil {
		state = int(DataUnknown)
	}
	m.State = State(state)
	return nil
}
