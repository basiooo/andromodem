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
	GetNetworkInfo(http.ResponseWriter, *http.Request)
	ToggleMobileData(http.ResponseWriter, *http.Request)
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
		util.MakeDeviceNotFoundResponse(writter)
		return
	}
	response := model.BaseResponse{
		Status:  "Success",
		Message: "Airplane mode status retrieved successfully",
		Data:    airplaneModeStatus,
	}
	util.WriteToResponseBody(writter, response, http.StatusOK)
}

func (d *NetworkHandlerImpl) ToggleAirplaneMode(writter http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	airplaneMode, err := d.NetworkService.ToggleAirplaneMode(serial)
	message := "Success disable airplane mode!"
	response := model.BaseResponse{
		Status:  "Success",
		Message: message,
	}
	if err != nil {
		if errors.Is(err, util.ErrDeviceNotFound) {
			util.MakeDeviceNotFoundResponse(writter)
			return
		} else {
			response.Status = "Failed"
			response.Message = err.Error()
		}
	}
	if airplaneMode.Enabled {
		message = "Success enable airplane mode!"
	}
	response.Data = airplaneMode
	util.WriteToResponseBody(writter, response, http.StatusOK)
}

func (d *NetworkHandlerImpl) GetNetworkInfo(writter http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	networkInfo, err := d.NetworkService.GetNetworkInfo(serial)
	if err != nil && errors.Is(err, util.ErrDeviceNotFound) {
		util.MakeDeviceNotFoundResponse(writter)
		return
	}
	response := model.BaseResponse{
		Status:  "Success",
		Message: "Network information retrieved successfully",
		Data: model.NetworkInfoResponse{
			NetworkInfo: *networkInfo,
		},
	}
	util.WriteToResponseBody(writter, response, http.StatusOK)
}

func (d *NetworkHandlerImpl) ToggleMobileData(writter http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	mobileData, err := d.NetworkService.ToggleMobileData(serial)

	message := "Success disable mobile data!"
	response := model.BaseResponse{
		Status:  "Success",
		Message: message,
		Data:    mobileData,
	}
	if err != nil {
		if errors.Is(err, util.ErrDeviceNotFound) {
			util.MakeDeviceNotFoundResponse(writter)
			return
		} else {
			response.Status = "Failed"
			response.Message = err.Error()
		}
	}
	if mobileData.Enabled {
		message = "Success enable mobile data!"
	}
	util.WriteToResponseBody(writter, response, http.StatusOK)
}
