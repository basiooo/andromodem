package router

import (
	"context"

	RESTHandler "github.com/basiooo/andromodem/internal/handler/rest"
	"github.com/basiooo/andromodem/internal/handler/web"
	"github.com/basiooo/andromodem/internal/service/devices_service"
	"github.com/basiooo/andromodem/internal/service/messages_service"
	"github.com/basiooo/andromodem/internal/service/monitoring_service"
	network_service "github.com/basiooo/andromodem/internal/service/network"
	"github.com/basiooo/andromodem/templates"
	"github.com/go-playground/validator/v10"

	_ "net/http/pprof"

	SSEHandler "github.com/basiooo/andromodem/internal/handler/sse"
	appMiddleware "github.com/basiooo/andromodem/internal/middleware"
	"github.com/basiooo/andromodem/pkg/adb_processor/processor"
	adb "github.com/basiooo/goadb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type Router struct {
	Adb       *adb.Adb
	Logger    *zap.Logger
	Ctx       context.Context
	ChiRouter chi.Router
	Validator *validator.Validate
}

func NewRouter(adb *adb.Adb, logger *zap.Logger, ctx context.Context, validator *validator.Validate) IRouter {
	return &Router{
		Adb:       adb,
		Logger:    logger,
		Ctx:       ctx,
		ChiRouter: chi.NewRouter(),
		Validator: validator,
	}
}

func (r *Router) GetRouters() chi.Router {
	r.ChiRouter.Use(
		middleware.Logger,
		// middleware.StripSlashes,
		appMiddleware.Recoverer(r.Logger),
		appMiddleware.AdbChecker(r.Adb, r.Logger),
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		}),
	)
	r.ChiRouter.Mount("/debug", middleware.Profiler())

	adbProcessor := processor.NewProcessor(r.Logger)

	// Services
	devicesService := devices_service.NewDevicesService(r.Adb, adbProcessor, r.Logger, r.Ctx)
	messagesService := messages_service.NewMessagesService(r.Adb, adbProcessor, r.Logger, r.Ctx)
	networkService := network_service.NewNetworkService(r.Adb, adbProcessor, r.Logger, r.Ctx)
	monitoringService := monitoring_service.NewMonitoringService(r.Adb, adbProcessor, networkService, r.Logger, r.Ctx)

	// Handlers
	devicesEventHandler := SSEHandler.NewDevicesEventHandler(devicesService, r.Logger)
	monitoringLogEventHandler := SSEHandler.NewMonitoringLogEventHandler(monitoringService, r.Logger)
	devicesHandler := RESTHandler.NewDevicesHandler(devicesService, r.Logger, r.Validator)
	messagesHandler := RESTHandler.NewMessagesHandler(messagesService, r.Logger, r.Validator)
	networkHandler := RESTHandler.NewNetworkHandler(networkService, r.Logger, r.Validator)
	monitoringHandler := RESTHandler.NewMonitoringHandler(monitoringService, r.Logger, r.Validator)

	// Frontend Handler
	frontendHandler := web.NewFrontendHandler(r.Logger, templates.MainPage)

	r.ChiRouter.Route("/api", func(r chi.Router) {
		r.Get("/devices/{serial}", devicesHandler.GetDeviceInfo)
		r.Post("/devices/{serial}/power", devicesHandler.PowerAction)
		r.Get("/devices/{serial}/feature-availabilities", devicesHandler.GetDeviceFeatureAvailabilities)
		r.Get("/devices/{serial}/messages", messagesHandler.GetMessages)
		r.Get("/devices/{serial}/network", networkHandler.GetNetworkInfo)
		r.Post("/devices/{serial}/network/mobile-data", networkHandler.ToggleMobileData)
		r.Post("/devices/{serial}/network/airplane-mode", networkHandler.ToggleAirplaneMode)

		r.Route("/devices/{serial}/monitoring", func(r chi.Router) {
			r.Post("/", monitoringHandler.CreateMonitoring)
			r.Get("/", monitoringHandler.GetMonitoringConfig)
			r.Put("/", monitoringHandler.UpdateMonitoringConfig)
			r.Post("/start", monitoringHandler.StartMonitoring)
			r.Post("/stop", monitoringHandler.StopMonitoring)
			r.Get("/status", monitoringHandler.GetMonitoringStatus)
			r.Get("/logs", monitoringHandler.GetMonitoringLogs)
			r.Delete("/logs", monitoringHandler.ClearMonitoringLogs)
		})

		r.Get("/monitoring", monitoringHandler.GetAllMonitoringTasks)
	})

	r.ChiRouter.Route("/event", func(r chi.Router) {
		r.Get("/devices", devicesEventHandler.ListenDevicesEvent)
		r.Route("/devices/{serial}/monitoring", func(r chi.Router) {
			r.Get("/logs", monitoringLogEventHandler.ListenMonitoringLogEvent)
		})
	})

	r.ChiRouter.Get("/assets/*", frontendHandler.ServeAssets().ServeHTTP)
	r.ChiRouter.Get("/", frontendHandler.ServeIndex)
	return r.ChiRouter
}
