package rest

import "net/http"

type IDevicesHandler interface {
	//GetDevices(http.ResponseWriter, *http.Request)
	GetDeviceInfo(http.ResponseWriter, *http.Request)
	PowerAction(http.ResponseWriter, *http.Request)
	GetDeviceFeatureAvailabilities(http.ResponseWriter, *http.Request)
}
