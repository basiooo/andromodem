package adb

import adb "github.com/abccyz/goadb"

type DeviceState adb.DeviceState

func (s DeviceState) String() string {
	switch adb.DeviceState(s) {
	case adb.StateDisconnected:
		return "Disconnect"
	case adb.StateOffline:
		return "Offline"
	case adb.StateOnline:
		return "Online"
	case adb.StateUnauthorized:
		return "Unauthorized"
	case adb.StateAuthorizing:
		return "Authorizing"
	default:
		return "Unknown"
	}
}
