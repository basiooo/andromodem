package monitoring_service

import (
	"context"

	"github.com/basiooo/andromodem/internal/model"
)

type IMonitoringWorkerService interface {
	StartMonitoring(string) error
	StopMonitoring(string, bool) error
	GetStatus(string) (*model.MonitoringStatus, error)
	IsRunning(string) bool
	AutoStartTasks() error
	Shutdown(context.Context) error
}
