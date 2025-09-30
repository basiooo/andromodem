package ws

import "net/http"

type IMirroringHandler interface {
	StartMirroringStream(http.ResponseWriter, *http.Request)
}
