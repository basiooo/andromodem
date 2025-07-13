package sse

import "net/http"

type IDevicesEventHandler interface {
	ListenDevicesEvent(http.ResponseWriter, *http.Request)
}
