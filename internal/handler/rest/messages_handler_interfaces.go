package rest

import "net/http"

type IMessagesHandler interface {
	GetMessages(http.ResponseWriter, *http.Request)
}
