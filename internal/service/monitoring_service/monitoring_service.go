package monitoring_service

import (
	"context"

	andromodemError "github.com/basiooo/andromodem/internal/errors"
	"github.com/basiooo/andromodem/internal/model"
	network_service "github.com/basiooo/andromodem/internal/service/network"
	"github.com/basiooo/andromodem/pkg/adb_processor/processor"
	adb "github.com/basiooo/goadb"
	"go.uber.org/zap"
)

type MonitoringService struct {
	taskService   IMonitoringTaskService
	configService IMonitoringConfigService
	workerService IMonitoringWorkerService
	logService    IMonitoringLogService
	logger        *zap.Logger
}

func NewMonitoringService(
	adb *adb.Adb,
	adbProcessor processor.IProcessor,
	networkService network_service.INetworkService,
	logger *zap.Logger,
	ctx context.Context,
) IMonitoringService {

	taskService := NewMonitoringTaskService(logger)
	configService := NewMonitoringConfigService("andromodem_monitoring_config.json", logger, taskService)
	logService := NewMonitoringLogService("andromodem_logs/monitoring", logger)
	pinggerService := NewMonitoringPinggerService(adb, logger)
	actionService := NewMonitoringDeviceActionService(adb, networkService, logService, logger)
	workerService := NewMonitoringWorkerService(ctx, logger, taskService, pinggerService, logService, actionService, configService)

	service := &MonitoringService{
		taskService:   taskService,
		configService: configService,
		workerService: workerService,
		logService:    logService,
		logger:        logger,
	}

	if err := service.LoadTasksFromFile(); err != nil {
		logger.Error("Failed to load tasks from config file", zap.Error(err))
	}

	if err := workerService.AutoStartTasks(); err != nil {
		logger.Error("Failed to auto start tasks", zap.Error(err))
	}

	return service
}

func (s *MonitoringService) CreateMonitoring(task *model.MonitoringTask) (*model.MonitoringTask, error) {
	createdTask, err := s.taskService.CreateTask(task)
	if err != nil {
		return nil, err
	}

	if err := s.SaveTasksToFile(); err != nil {
		s.logger.Error("[Monitoring] Failed to save tasks to file", zap.Error(err))
		return nil, err
	}

	if createdTask.IsActive {
		if err := s.workerService.StartMonitoring(createdTask.Serial); err != nil {
			s.logger.Error("[Monitoring] Failed to auto-start monitoring after creation",
				zap.String("serial", createdTask.Serial), zap.Error(err))
		}
	}

	return createdTask, nil
}

func (s *MonitoringService) StartMonitoring(serial string) error {
	if err := s.workerService.StartMonitoring(serial); err != nil {
		return err
	}

	if err := s.SaveTasksToFile(); err != nil {
		s.logger.Error("[Monitoring] Failed to save tasks to file after start", zap.Error(err))
	}

	return nil
}

func (s *MonitoringService) StopMonitoring(serial string) error {

	if err := s.workerService.StopMonitoring(serial, false); err != nil {
		return err
	}

	if err := s.SaveTasksToFile(); err != nil {
		s.logger.Error("[Monitoring] Failed to save tasks to file after stop", zap.Error(err))
	}

	return nil
}

func (s *MonitoringService) GetMonitoringStatus(serial string) (*model.MonitoringStatus, error) {
	return s.workerService.GetStatus(serial)
}

func (s *MonitoringService) DeleteMonitoring(serial string) error {
	if err := s.StopMonitoring(serial); err != nil {
		s.logger.Warn("[Monitoring] Failed to stop monitoring task during deletion",
			zap.String("serial", serial), zap.Error(err))
	}

	if err := s.taskService.DeleteTask(serial); err != nil {
		return err
	}

	return s.SaveTasksToFile()
}

func (s *MonitoringService) GetMonitoringConfig(serial string) (*model.MonitoringTask, error) {
	return s.taskService.GetTask(serial)
}

func (s *MonitoringService) UpdateMonitoringConfig(serial string, request *model.MonitoringTaskRequest) (*model.MonitoringTask, error) {
	wasRunning := s.workerService.IsRunning(serial)

	if wasRunning {
		if err := s.workerService.StopMonitoring(serial, true); err != nil {
			s.logger.Warn("[Monitoring] Failed to stop monitoring during update", zap.Error(err))
		}
	}

	task, err := s.taskService.UpdateTask(serial, request)
	if err != nil {
		return nil, err
	}

	if err := s.SaveTasksToFile(); err != nil {
		s.logger.Error("[Monitoring] Failed to save tasks to file after update", zap.Error(err))
	}

	if wasRunning {
		if err := s.workerService.StartMonitoring(serial); err != nil {
			s.logger.Error("[Monitoring] Failed to restart monitoring after update", zap.Error(err))
		}
	}

	return task, nil
}

func (s *MonitoringService) GetAllMonitoringTasks() ([]*model.MonitoringTask, error) {
	return s.taskService.GetAllTasks()
}

func (s *MonitoringService) GetMonitoringLogs(serial string, limit int) ([]*model.MonitoringLog, error) {
	return s.logService.GetLogs(serial, limit)
}

func (s *MonitoringService) ListenMonitoringLogs(ctx context.Context, serial string, callback func(*model.MonitoringLog) error) error {
	return s.logService.LogListener(ctx, serial, callback)
}

func (s *MonitoringService) LoadTasksFromFile() error {
	tasks, err := s.configService.LoadTasksFromFile()
	if err != nil {
		return err
	}

	if taskServiceImpl, ok := s.taskService.(*MonitoringTaskService); ok {
		taskServiceImpl.LoadTasks(tasks)
	}

	return nil
}

func (s *MonitoringService) SaveTasksToFile() error {
	tasks, err := s.taskService.GetAllTasks()
	if err != nil {
		return err
	}

	return s.configService.SaveTasksToFile(tasks)
}

func (s *MonitoringService) ClearMonitoringLogs(serial string) error {
	_, err := s.taskService.GetTask(serial)
	if err != nil {
		return andromodemError.ErrorTaskNotFoundInConfig
	}
	return s.logService.ClearLogs(serial)
}
