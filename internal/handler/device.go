package handler

import (
	"errors"
	"net/http"

	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/service"
	"github.com/basiooo/andromodem/internal/util"
	"github.com/go-chi/chi/v5"
)

type DeviceHandler interface {
	GetDevices(http.ResponseWriter, *http.Request)
	GetDeviceInfo(http.ResponseWriter, *http.Request)
}

type DeviceHandlerImpl struct {
	DeviceService service.DeviceService
}

func NewDeviceHander(deviceService service.DeviceService) DeviceHandler {
	return &DeviceHandlerImpl{
		DeviceService: deviceService,
	}
}

func (d *DeviceHandlerImpl) GetDevices(writter http.ResponseWriter, request *http.Request) {
	devices, _ := d.DeviceService.GetDevices()
	response := model.BaseResponse{
		Status:  "Success",
		Message: "Device list retrieved successfully",
		Data: model.DevicesResponse{
			Devices: devices,
		},
	}
	util.WriteToResponseBody(writter, response, http.StatusOK)
}

func (d *DeviceHandlerImpl) GetDeviceInfo(writter http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	deviceInfo, err := d.DeviceService.GetDeviceInfo(serial)
	if err != nil && errors.Is(err, util.ErrDeviceNotFound) {
		util.MakeDeviceNotFoundResponse(writter)
		return
	}
	response := model.BaseResponse{
		Status:  "Success",
		Message: "Device information retrieved successfully",
		Data: model.DeviceInfoResponse{
			DeviceInfo: *deviceInfo,
		},
	}
	util.WriteToResponseBody(writter, response, http.StatusOK)
}
