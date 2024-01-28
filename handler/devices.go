package handler

import (
	adb "github.com/abccyz/goadb"
	"github.com/basiooo/andromodem/helper"
	"github.com/basiooo/andromodem/model"
	"github.com/sirupsen/logrus"
	"net/http"
)

type DeviceState struct {
	adb.DeviceState
}

func (s DeviceState) String() string {
	switch s.DeviceState {
	case adb.StateDisconnected:
		return "disconnect"
	case adb.StateOffline:
		return "offline"
	case adb.StateOnline:
		return "online"
	case adb.StateUnauthorized:
		return "unauthorized"
	case adb.StateAuthorizing:
		return "authorizing"
	default:
		return "unknown"
	}
}

func GetDevices(writter http.ResponseWriter, request *http.Request) {
	adbClient, err := helper.GetADBClient(request)
	helper.PanicIfError(err)
	deviceListInfo, err := adbClient.ListDevices()
	if err != nil {
		logrus.Errorf("GetDevices(): failed get list devices. error: %v", err)
		helper.PanicIfError(err)
	}

	devices := make([]model.Devices, 0, len(deviceListInfo))
	for _, device := range deviceListInfo {
		deviceInfo := adbClient.Device(adb.DeviceWithSerial(device.Serial))
		_deviceState, err := deviceInfo.State()
		if err != nil {
			logrus.WithField("function", "GetDevices").Errorf("GetDevices(): failed get device state %v. error: %v", device.Model, err)
			_deviceState = adb.StateInvalid
		}
		deviceState := DeviceState{
			_deviceState,
		}
		deviceResponse := model.Devices{
			Serial: device.Serial,
			Model:  device.Model,
			State:  deviceState.String(),
		}
		devices = append(devices, deviceResponse)
	}
	response := model.DevicesResponse{
		Devices: devices,
	}
	helper.WriteToResponseBody(writter, response, 200)
}
