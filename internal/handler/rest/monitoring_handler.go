package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/basiooo/andromodem/internal/common"
	andromodemError "github.com/basiooo/andromodem/internal/errors"
	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/service/monitoring_service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type MonitoringHandler struct {
	MonitoringService monitoring_service.IMonitoringService
	Logger            *zap.Logger
	Validator         *validator.Validate
}

func NewMonitoringHandler(monitoringService monitoring_service.IMonitoringService, logger *zap.Logger, validator *validator.Validate) IMonitoringHandler {
	return &MonitoringHandler{
		MonitoringService: monitoringService,
		Logger:            logger,
		Validator:         validator,
	}
}

func (h *MonitoringHandler) StartMonitoring(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")

	if err := h.MonitoringService.StartMonitoring(serial); err != nil {

		if errors.Is(err, andromodemError.ErrorTaskNotFoundInConfig) {
			common.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		} else if errors.Is(err, andromodemError.ErrorMonitoringTaskAlreadyRunning) {
			common.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Logger.Error("failed to start monitoring", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.SuccessResponse(w, "Monitoring started successfully", nil, http.StatusOK)
}

func (h *MonitoringHandler) StopMonitoring(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")

	if err := h.MonitoringService.StopMonitoring(serial); err != nil {
		if errors.Is(err, andromodemError.ErrorTaskNotFoundInConfig) {
			common.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		} else if errors.Is(err, andromodemError.ErrorNoRunningMonitoringTask) {
			common.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Logger.Error("failed to stop monitoring", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.SuccessResponse(w, "Monitoring stopped successfully", nil, http.StatusOK)
}

func (h *MonitoringHandler) DeleteMonitoring(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")

	if err := h.MonitoringService.DeleteMonitoring(serial); err != nil {
		if errors.Is(err, andromodemError.ErrorTaskNotFoundInConfig) {
			common.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		h.Logger.Error("failed to delete monitoring", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.SuccessResponse(w, "Monitoring deleted successfully", nil, http.StatusOK)
}

func (h *MonitoringHandler) GetMonitoringStatus(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")

	status, err := h.MonitoringService.GetMonitoringStatus(serial)
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.DeviceNotFoundResponse(w)
			return
		}
		if errors.Is(err, andromodemError.ErrorTaskNotFoundInConfig) {
			common.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		h.Logger.Error("failed to get monitoring status", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.SuccessResponse(w, "Monitoring status retrieved successfully", status, http.StatusOK)
}

func (h *MonitoringHandler) GetMonitoringConfig(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")

	task, err := h.MonitoringService.GetMonitoringConfig(serial)
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.ErrorResponse(w, "Monitoring config not found for this device", http.StatusNotFound)
			return
		}

		if errors.Is(err, andromodemError.ErrorTaskNotFoundInConfig) {
			common.ErrorResponse(w, err.Error(), http.StatusNotFound)
			return
		}
		h.Logger.Error("failed to get monitoring config", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.SuccessResponse(w, "Monitoring config retrieved successfully", task, http.StatusOK)
}

func (h *MonitoringHandler) UpdateMonitoringConfig(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")

	var request model.MonitoringTaskRequest
	err := common.ReadFromRequestBody(r, &request)
	if err != nil {
		h.Logger.Error("failed to read request body", zap.Error(err))
		common.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.Validator.Struct(request); err != nil {
		h.Logger.Error("validation error", zap.Error(err))
		common.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTask, err := h.MonitoringService.UpdateMonitoringConfig(serial, &request)
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.ErrorResponse(w, "Monitoring config not found for this device", http.StatusNotFound)
			return
		}
		h.Logger.Error("failed to update monitoring config", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.SuccessResponse(w, "Monitoring config updated successfully", updatedTask, http.StatusOK)
}

func (h *MonitoringHandler) GetAllMonitoringTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.MonitoringService.GetAllMonitoringTasks()
	if err != nil {
		h.Logger.Error("failed to get all monitoring tasks", zap.Error(err))
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.SuccessResponse(w, "Monitoring tasks retrieved successfully", tasks, http.StatusOK)
}

func (h *MonitoringHandler) GetMonitoringLogs(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")
	if serial == "" {
		common.ErrorResponse(w, "Serial is required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := h.MonitoringService.GetMonitoringLogs(serial, limit)
	if err != nil {
		h.Logger.Error("Failed to get monitoring logs", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(w, "Failed to get monitoring logs", http.StatusInternalServerError)
		return
	}

	common.SuccessResponse(w, "Monitoring logs retrieved successfully", logs, http.StatusOK)
}

func (h *MonitoringHandler) CreateMonitoring(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")

	var request model.MonitoringTaskRequest
	err := common.ReadFromRequestBody(r, &request)
	if err != nil {
		h.Logger.Error("failed to read request body", zap.Error(err))
		common.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Validator.Struct(request); err != nil {
		h.Logger.Error("validation error", zap.Error(err))
		common.ValidationErrorResponse(w, "Validation error", http.StatusBadRequest, err)
		return
	}

	task := &model.MonitoringTask{
		Serial:            serial,
		Host:              request.Host,
		Method:            request.Method,
		MaxFailures:       request.MaxFailures,
		CheckingInterval:  request.CheckingInterval,
		AirplaneModeDelay: request.AirplaneModeDelay,
	}

	createdTask, err := h.MonitoringService.CreateMonitoring(task)
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.ErrorResponse(w, "Device not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, andromodemError.ErrorMonitoringTaskExists) {
			common.ErrorResponse(w, "Monitoring config already exist", http.StatusBadRequest)
			return
		}
		h.Logger.Error("failed to create monitoring", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.SuccessResponse(w, "Monitoring task created successfully", createdTask, http.StatusCreated)
}

func (h *MonitoringHandler) ClearMonitoringLogs(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")

	err := h.MonitoringService.ClearMonitoringLogs(serial)
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.ErrorResponse(w, "Monitoring config not found for this device", http.StatusNotFound)
			return
		}
		h.Logger.Error("failed to clear monitoring logs", zap.String("serial", serial), zap.Error(err))
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.SuccessResponse(w, "Monitoring logs cleared successfully", nil, http.StatusOK)
}
