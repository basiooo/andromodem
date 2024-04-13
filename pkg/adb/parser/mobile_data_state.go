package parser

import (
	"strconv"

	"github.com/sirupsen/logrus"
)

type ConnectionState int8

const (
	DataUnknown      ConnectionState = -1
	DataDisconnected ConnectionState = 0
	DataConnecting   ConnectionState = 1
	DataConnected    ConnectionState = 2
	DataSuspended    ConnectionState = 3
)

var connectionStateStrings = map[ConnectionState]string{
	DataUnknown:      "Unknown",
	DataDisconnected: "Disconnected",
	DataConnecting:   "Connecting",
	DataConnected:    "Connected",
	DataSuspended:    "Suspended",
}

func (connectionState ConnectionState) String() string {
	if val, ok := connectionStateStrings[connectionState]; ok {
		return val
	}
	return "Unknown"
}

type MobileDataConnectionState struct {
	ConnectionState ConnectionState
}

func NewMobileDataConnectionState(rawState string) *MobileDataConnectionState {
	mDataConnectionState := &MobileDataConnectionState{
		ConnectionState: parserState(rawState),
	}
	return mDataConnectionState
}

func parserState(rawState string) ConnectionState {
	state, err := strconv.Atoi(rawState)
	if err != nil {
		logrus.WithField("location", "mobile_data_state.parseState").Errorf("cannot convert %v to int: %v", rawState, err)
		state = int(DataUnknown)
	}
	return ConnectionState(state)
}
