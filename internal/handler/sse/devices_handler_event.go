package sse

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/basiooo/andromodem/internal/common"
	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/service/devices_service"
	"go.uber.org/zap"
)

type DevicesEventHandler struct {
	DevicesEventService devices_service.IDevicesService
	Logger              *zap.Logger
}

func NewDevicesEventHandler(devicesEventService devices_service.IDevicesService, logger *zap.Logger) IDevicesEventHandler {
	return &DevicesEventHandler{
		DevicesEventService: devicesEventService,
		Logger:              logger,
	}
}

func (d *DevicesEventHandler) ListenDevicesEvent(w http.ResponseWriter, r *http.Request) {
	common.SSESetResponseHeader(w)
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusServiceUnavailable)
		return
	}
	requestCtx := r.Context()
	_, err := fmt.Fprintf(w, "data: connected\n\n")
	if err != nil {
		d.Logger.Error("error writing to response:", zap.Error(err))
		return
	}
	flusher.Flush()
	err = d.DevicesEventService.DevicesListener(requestCtx, func(device *model.Device) error {
		res, err := json.Marshal(device)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "data: %s\n\n", res)
		if err != nil {
			d.Logger.Error("error writing to response:", zap.Error(err))
			return err
		}
		flusher.Flush()
		return err
	})

	if err != nil {
		d.Logger.Error("Listen stopped:", zap.Error(err))
	}

}
