package parser

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type MobileDataConnectionState int8

const (
	DataUnknown      MobileDataConnectionState = -1
	DataDisconnected MobileDataConnectionState = 0
	DataConnecting   MobileDataConnectionState = 1
	DataConnected    MobileDataConnectionState = 2
	DataSuspended    MobileDataConnectionState = 3
)

var mobileDataConnectionStateStrings = map[MobileDataConnectionState]string{
	DataUnknown:      "Unknown",
	DataDisconnected: "Disconnected",
	DataConnecting:   "Connecting",
	DataConnected:    "Connected",
	DataSuspended:    "Suspended",
}

func (connectionState MobileDataConnectionState) String() string {
	if val, ok := mobileDataConnectionStateStrings[connectionState]; ok {
		return val
	}
	return "Unknown"
}

func ParseMobileDataConnectionState(rawMobileDataConnectionState string) []MobileDataConnectionState {
	var mobileDataConnectionState []MobileDataConnectionState
	for _, data := range strings.Split(strings.TrimSpace(rawMobileDataConnectionState), "\n") {
		state, err := strconv.Atoi(strings.TrimSpace(data))
		if err != nil {
			logrus.WithField("function", "ParseMobileDataConnectionState").Errorf("cannot convert %v to int: %v", data, err)
			state = int(DataUnknown)
		}
		mobileDataConnectionState = append(mobileDataConnectionState, MobileDataConnectionState(state))
	}
	return mobileDataConnectionState
}
