package router

import "github.com/go-chi/chi/v5"

type IRouter interface {
	GetRouters() chi.Router
}
