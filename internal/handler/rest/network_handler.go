package rest

import (
	"errors"
	"net/http"

	"github.com/basiooo/andromodem/internal/common"
	andromodemError "github.com/basiooo/andromodem/internal/errors"
	network_service "github.com/basiooo/andromodem/internal/service/network"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type NetworkHandler struct {
	NetworkService network_service.INetworkService
	Logger         *zap.Logger
	Validator      *validator.Validate
}

func NewNetworkHandler(networkService network_service.INetworkService, logger *zap.Logger, validator *validator.Validate) INetworkHandler {
	return &NetworkHandler{
		NetworkService: networkService,
		Logger:         logger,
		Validator:      validator,
	}
}

func (n *NetworkHandler) GetNetworkInfo(writer http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	networkInfo, err := n.NetworkService.GetNetworkInfo(serial)
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.DeviceNotFoundResponse(writer)
			return
		}
		n.Logger.Error("error getting network info", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(writer, "Error getting network info", http.StatusInternalServerError)
		return
	}
	common.SuccessResponse(writer, "Network info retrieved successfully", networkInfo, http.StatusOK)
}

func (n *NetworkHandler) ToggleMobileData(writer http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	toggleResult, err := n.NetworkService.ToggleMobileData(serial)
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.DeviceNotFoundResponse(writer)
			return
		}
		n.Logger.Error("error toggling mobile data", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	message := "Mobile data disabled successfully"
	if *toggleResult {
		message = "Mobile data enabled successfully"
	}
	common.SuccessResponse(writer, message, nil, http.StatusOK)
}

func (n *NetworkHandler) ToggleAirplaneMode(writer http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	toggleResult, err := n.NetworkService.ToggleAirplaneMode(serial)
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.DeviceNotFoundResponse(writer)
			return
		}
		n.Logger.Error("error toggling airplane mode", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	message := "Airplane mode disabled successfully"
	if *toggleResult {
		message = "Airplane mode enabled successfully"
	}
	common.SuccessResponse(writer, message, nil, http.StatusOK)
}
