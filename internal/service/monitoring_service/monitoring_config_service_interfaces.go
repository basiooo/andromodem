package monitoring_service

import "github.com/basiooo/andromodem/internal/model"

type IMonitoringConfigService interface {
	LoadTasksFromFile() ([]*model.MonitoringTask, error)
	SaveTasksToFile([]*model.MonitoringTask) error
	ValidateAndCleanConfig() error
	SetConfigFile(string)
	GetConfigFile() string
}
