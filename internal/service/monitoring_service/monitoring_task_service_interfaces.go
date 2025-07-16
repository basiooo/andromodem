package monitoring_service

import "github.com/basiooo/andromodem/internal/model"

type IMonitoringTaskService interface {
	CreateTask(*model.MonitoringTask) (*model.MonitoringTask, error)
	UpdateTask(string, *model.MonitoringTaskRequest) (*model.MonitoringTask, error)
	UpdateTaskStatus(string, bool) error
	UpdateTaskField(string, func(*model.MonitoringTask)) error
	DeleteTask(string) error
	GetTask(string) (*model.MonitoringTask, error)
	GetAllTasks() ([]*model.MonitoringTask, error)
	ValidateTask(*model.MonitoringTask) bool
	TaskExists(string) bool
	LoadTasks(tasks []*model.MonitoringTask)
}
