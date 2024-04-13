package server

import (
	"net/http"

	"github.com/basiooo/andromodem/internal/adb"
	"github.com/basiooo/andromodem/internal/handler"
	appMiddleware "github.com/basiooo/andromodem/internal/middleware"
	"github.com/basiooo/andromodem/internal/service"
	adbcommand "github.com/basiooo/andromodem/pkg/adb/adb_command"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Router interface {
	Setup() *chi.Mux
}
type routerImpl struct {
	*adb.Adb
}

func NewRouter(adbClient *adb.Adb) Router {
	return &routerImpl{
		Adb: adbClient,
	}
}
func (r *routerImpl) Setup() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.StripSlashes,
		appMiddleware.Recoverer,
	)
	adbcommand := adbcommand.NewAdbCommand()
	devicesService := service.NewDeviceService(r.Adb, *adbcommand)
	devicesHandler := handler.NewDeviceHander(devicesService)

	router.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})
		r.Get("/devices", devicesHandler.GetDevices)
		r.Get("/devices/{serial}", devicesHandler.GetDeviceInfo)
	})
	return router
}
