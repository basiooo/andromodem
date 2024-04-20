package server

import (
	"embed"
	"net/http"
	"text/template"

	"github.com/basiooo/andromodem/internal/adb"
	"github.com/basiooo/andromodem/internal/handler"
	appMiddleware "github.com/basiooo/andromodem/internal/middleware"
	"github.com/basiooo/andromodem/internal/service"
	adbcommand "github.com/basiooo/andromodem/pkg/adb/adb_command"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

type Router interface {
	Setup() *chi.Mux
}
type routerImpl struct {
	*adb.Adb
	TemplateFs embed.FS
}

func NewRouter(adb *adb.Adb, templateFS embed.FS) Router {
	return &routerImpl{
		Adb:        adb,
		TemplateFs: templateFS,
	}
}
func (r *routerImpl) Setup() *chi.Mux {
	adb := r.Adb
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.StripSlashes,
		appMiddleware.Recoverer,
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		}),
	)
	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		t, err := template.ParseFS(r.TemplateFs, "andromodem-dashboard/dist/index.html")
		if err != nil {
			logrus.WithField("function", "NewRouter").Fatal("failed get template : ", err)
		}
		err = t.Execute(writer, "")
		if err != nil {
			logrus.WithField("function", "NewRouter").Fatal("failed execute template : ", err)
		}
	})
	router.Route("/api", func(r chi.Router) {
		r.Use(appMiddleware.AdbChecker(adb))
		adbcommand := adbcommand.NewAdbCommand()
		devicesService := service.NewDeviceService(adb, *adbcommand)
		devicesHandler := handler.NewDeviceHander(devicesService)
		messageService := service.NewMessageService(adb, *adbcommand)
		messageHandler := handler.NewMessageHander(messageService)
		networkService := service.NewNetworkService(adb, *adbcommand)
		networkHandler := handler.NewNetworkHander(networkService)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK")) //nolint:errcheck
		})
		r.Get("/devices", devicesHandler.GetDevices)
		r.Get("/devices/{serial}", devicesHandler.GetDeviceInfo)
		r.Get("/devices/{serial}/inbox", messageHandler.GetSmsInbox)
		r.Get("/devices/{serial}/network/airplane", networkHandler.GetAirplaneModeStatus)
		r.Put("/devices/{serial}/network/airplane", networkHandler.ToggleAirplaneMode)
		r.Get("/devices/{serial}/network", networkHandler.GetNetworkInfo)
		r.Put("/devices/{serial}/network", networkHandler.ToggleMobileData)
	})
	return router
}
