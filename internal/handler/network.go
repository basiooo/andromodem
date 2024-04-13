package handler

import (
	"errors"
	"net/http"

	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/service"
	"github.com/basiooo/andromodem/internal/util"
	"github.com/go-chi/chi/v5"
)

type NetworkHandler interface {
	GetAirplaneModeStatus(http.ResponseWriter, *http.Request)
	ToggleAirplaneMode(http.ResponseWriter, *http.Request)
}

type NetworkHandlerImpl struct {
	NetworkService service.NetworkService
}

func NewNetworkHander(deviceService service.NetworkService) NetworkHandler {
	return &NetworkHandlerImpl{
		NetworkService: deviceService,
	}
}

func (d *NetworkHandlerImpl) GetAirplaneModeStatus(writter http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	airplaneModeStatus, err := d.NetworkService.GetAirplaneModeStatus(serial)
	if err != nil && errors.Is(err, util.ErrDeviceNotFound) {
		response := model.ErrorResponse{
			Error: "Device Not Found",
		}
		util.WriteToResponseBody(writter, response, http.StatusNotFound)
		return
	}
	util.WriteToResponseBody(writter, airplaneModeStatus, http.StatusOK)
}

func (d *NetworkHandlerImpl) ToggleAirplaneMode(writter http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	airplaneModeStatus, err := d.NetworkService.ToggleAirplaneMode(serial)
	if err != nil && errors.Is(err, util.ErrDeviceNotFound) {
		response := model.ErrorResponse{
			Error: "Device Not Found",
		}
		util.WriteToResponseBody(writter, response, http.StatusNotFound)
		return
	}
	util.WriteToResponseBody(writter, airplaneModeStatus, http.StatusOK)
}
