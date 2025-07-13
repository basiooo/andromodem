package monitoring_service

import (
	"context"

	"github.com/basiooo/andromodem/internal/model"
)

type IMonitoringService interface {
	CreateMonitoring(*model.MonitoringTask) (*model.MonitoringTask, error)
	StartMonitoring(string) error
	StopMonitoring(string) error
	DeleteMonitoring(string) error
	ClearMonitoringLogs(string) error
	GetMonitoringStatus(string) (*model.MonitoringStatus, error)
	GetMonitoringConfig(string) (*model.MonitoringTask, error)
	UpdateMonitoringConfig(string, *model.MonitoringTaskRequest) (*model.MonitoringTask, error)
	GetAllMonitoringTasks() ([]*model.MonitoringTask, error)
	GetMonitoringLogs(string, int) ([]*model.MonitoringLog, error)
	LoadTasksFromFile() error
	SaveTasksToFile() error
	ListenMonitoringLogs(ctx context.Context, serial string, callback func(*model.MonitoringLog) error) error
}
