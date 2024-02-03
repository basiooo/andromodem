package app

import (
	"embed"
	adb "github.com/abccyz/goadb"
	"github.com/basiooo/andromodem/handler"
	cmiddleware "github.com/basiooo/andromodem/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

func NewRouter(templateFS embed.FS, adbClient *adb.Adb) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.StripSlashes, cmiddleware.Recoverer, cmiddleware.AdbClient(adbClient))
	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		t, err := template.ParseFS(templateFS, "andromodem-frontend/dist/index.html")
		if err != nil {
			logrus.WithField("function", "NewRouter").Fatal("failed get template : ", err)
		}
		err = t.Execute(writer, "")
		if err != nil {
			logrus.WithField("function", "NewRouter").Fatal("failed execute template : ", err)
		}
	})
	r.Route("/api", func(r chi.Router) {
		r.Get("/devices", handler.GetDevices)
		r.Get("/device/{serial}", handler.GetDevice)
		r.Get("/device/{serial}/thermal", handler.GetThermal)
		r.Get("/network/{serial}", handler.GetNetwork)
		r.Post("/network/{serial}/mobile-data/toggle", handler.MobileDataToggle)
	})
	return r
}
