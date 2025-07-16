package rest

import "net/http"

type INetworkHandler interface {
	GetNetworkInfo(http.ResponseWriter, *http.Request)
	ToggleMobileData(http.ResponseWriter, *http.Request)
	ToggleAirplaneMode(http.ResponseWriter, *http.Request)
}
