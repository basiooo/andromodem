package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	andromodemError "github.com/basiooo/andromodem/internal/errors"
	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/service/devices_service"
	"github.com/go-playground/validator/v10"

	"github.com/basiooo/andromodem/internal/common"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type DevicesHandler struct {
	DevicesService devices_service.IDevicesService
	Logger         *zap.Logger
	Validator      *validator.Validate
}

func NewDevicesHandler(devicesService devices_service.IDevicesService, logger *zap.Logger, validator *validator.Validate) IDevicesHandler {
	return &DevicesHandler{
		DevicesService: devicesService,
		Logger:         logger,
		Validator:      validator,
	}
}

func (d *DevicesHandler) GetDeviceInfo(writer http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	deviceInfo, err := d.DevicesService.GetDeviceInfo(serial)
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.DeviceNotFoundResponse(writer)
			return
		}
		d.Logger.Error("error getting device info", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(writer, "Error getting device info", http.StatusInternalServerError)
		return
	}
	common.SuccessResponse(writer, "Device info retrieved successfully", deviceInfo, http.StatusOK)
}

func (d *DevicesHandler) GetDeviceFeatureAvailabilities(writer http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	DeviceFeatuAvailabilities, err := d.DevicesService.GetDeviceFeatureAvailabilities(serial)
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.DeviceNotFoundResponse(writer)
			return
		}
		d.Logger.Error("error getting device feature availabilities", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(writer, "Error getting device feature availabilities", http.StatusInternalServerError)
		return
	}
	common.SuccessResponse(writer, "Device feature availabilities retrieved successfully", DeviceFeatuAvailabilities, http.StatusOK)
}

func (d *DevicesHandler) PowerAction(writer http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")

	var powerAction model.DevicePowerAction
	err := json.NewDecoder(request.Body).Decode(&powerAction)
	if err != nil {
		d.Logger.Error("error decoding power action request", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(writer, "Error decoding request", http.StatusBadRequest)
		return
	}
	if err := d.Validator.Struct(powerAction); err != nil {
		d.Logger.Error("error validating power action request", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(writer, "Error validating power action request", http.StatusBadRequest)
		return
	}
	err = d.DevicesService.DevicePower(serial, devices_service.PowerAction(powerAction.Action))
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.DeviceNotFoundResponse(writer)
			return
		}
		d.Logger.Error("error power action", zap.String("serial", serial), zap.String("action", string(powerAction.Action)), zap.Error(err))
		common.ErrorResponse(writer, "Error power action", http.StatusInternalServerError)
		return
	}
	common.SuccessResponse(writer, "Power action executed successfully", nil, http.StatusOK)
}
