package router

import (
	"context"

	"github.com/basiooo/andromodem/internal/handler/rest"
	"github.com/basiooo/andromodem/internal/handler/ws"

	"github.com/basiooo/andromodem/internal/handler/web"
	"github.com/basiooo/andromodem/internal/service/devices_service"
	"github.com/basiooo/andromodem/internal/service/messages_service"
	"github.com/basiooo/andromodem/internal/service/mirroring_service"
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
	mirroringService := mirroring_service.NewMirroringService(r.Adb, r.Logger, r.Ctx)

	// Handlers
	devicesEventHandler := SSEHandler.NewDevicesEventHandler(devicesService, r.Logger)
	monitoringLogEventHandler := SSEHandler.NewMonitoringLogEventHandler(monitoringService, r.Logger)
	devicesHandler := rest.NewDevicesHandler(devicesService, r.Logger, r.Validator)
	messagesHandler := rest.NewMessagesHandler(messagesService, r.Logger, r.Validator)
	networkHandler := rest.NewNetworkHandler(networkService, r.Logger, r.Validator)
	monitoringHandler := rest.NewMonitoringHandler(monitoringService, r.Logger, r.Validator)
	mirroringHandler := ws.NewMirroringHandler(mirroringService, r.Logger, r.Validator)

	healthHandler := rest.NewHealthHandler()

	// Frontend Handler
	frontendHandler := web.NewFrontendHandler(r.Logger, templates.MainPage)

	r.ChiRouter.Route("/api", func(chiRouter chi.Router) {
		chiRouter.Use(appMiddleware.AdbChecker(r.Adb, r.Logger))
		chiRouter.Route("/devices/{serial}", func(chiRouter chi.Router) {
			chiRouter.Get("/", devicesHandler.GetDeviceInfo)
			chiRouter.Post("/power", devicesHandler.PowerAction)
			chiRouter.Get("/feature-availabilities", devicesHandler.GetDeviceFeatureAvailabilities)
			chiRouter.Get("/messages", messagesHandler.GetMessages)
			chiRouter.Get("/network", networkHandler.GetNetworkInfo)
			chiRouter.Post("/network/mobile-data", networkHandler.ToggleMobileData)
			chiRouter.Post("/network/airplane-mode", networkHandler.ToggleAirplaneMode)

			chiRouter.Route("/monitoring", func(chiRouter chi.Router) {
				chiRouter.Post("/", monitoringHandler.CreateMonitoring)
				chiRouter.Get("/", monitoringHandler.GetMonitoringConfig)
				chiRouter.Put("/", monitoringHandler.UpdateMonitoringConfig)
				chiRouter.Post("/start", monitoringHandler.StartMonitoring)
				chiRouter.Post("/stop", monitoringHandler.StopMonitoring)
				chiRouter.Get("/status", monitoringHandler.GetMonitoringStatus)
				chiRouter.Get("/logs", monitoringHandler.GetMonitoringLogs)
				chiRouter.Delete("/logs", monitoringHandler.ClearMonitoringLogs)
			})
		})
		chiRouter.Get("/health/ping", healthHandler.Ping)
		chiRouter.Get("/monitoring", monitoringHandler.GetAllMonitoringTasks)
	})

	r.ChiRouter.Route("/event", func(chiRouter chi.Router) {
		chiRouter.Get("/devices", devicesEventHandler.ListenDevicesEvent)
		chiRouter.Route("/devices/{serial}/monitoring", func(chiRouter chi.Router) {
			chiRouter.Get("/logs", monitoringLogEventHandler.ListenMonitoringLogEvent)
		})
	})

	r.ChiRouter.Route("/ws", func(chiRouter chi.Router) {
		chiRouter.Route("/devices/{serial}", func(chiRouter chi.Router) {
			chiRouter.Get("/mirroring", mirroringHandler.StartMirroringStream)
		})
	})
	r.ChiRouter.Get("/assets/*", frontendHandler.ServeAssets().ServeHTTP)
	r.ChiRouter.Get("/", frontendHandler.ServeIndex)
	return r.ChiRouter
}
