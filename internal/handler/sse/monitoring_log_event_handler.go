package sse

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/basiooo/andromodem/internal/common"
	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/service/monitoring_service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type MonitoringLogEventHandler struct {
	MonitoringService monitoring_service.IMonitoringService
	Logger            *zap.Logger
}

func NewMonitoringLogEventHandler(monitoringService monitoring_service.IMonitoringService, logger *zap.Logger) IMonitoringLogEventHandler {
	return &MonitoringLogEventHandler{
		MonitoringService: monitoringService,
		Logger:            logger,
	}
}

func (h *MonitoringLogEventHandler) ListenMonitoringLogEvent(w http.ResponseWriter, r *http.Request) {
	common.SSESetResponseHeader(w)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusServiceUnavailable)
		return
	}
	serial := chi.URLParam(r, "serial")

	requestCtx := r.Context()

	err := h.MonitoringService.ListenMonitoringLogs(requestCtx, serial, func(log *model.MonitoringLog) error {

		res, err := json.Marshal(log)
		if err != nil {
			h.Logger.Error("error marshaling log:", zap.Error(err))
			return err
		}

		_, err = fmt.Fprintf(w, "data: %s\n\n", res)
		if err != nil {
			h.Logger.Error("error writing log to response:", zap.Error(err))
			return err
		}
		flusher.Flush()
		return nil
	})

	if err != nil {
		h.Logger.Error("Log listener stopped:", zap.Error(err))
	}
}
