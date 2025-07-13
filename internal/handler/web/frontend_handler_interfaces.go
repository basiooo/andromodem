package web

import "net/http"

type IFrontendHandler interface {
	ServeIndex(w http.ResponseWriter, r *http.Request)
	ServeAssets() http.Handler
}
