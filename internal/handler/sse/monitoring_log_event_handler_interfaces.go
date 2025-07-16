package sse

import "net/http"

type IMonitoringLogEventHandler interface {
	ListenMonitoringLogEvent(http.ResponseWriter, *http.Request)
}
