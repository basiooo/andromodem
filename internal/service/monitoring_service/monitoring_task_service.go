package monitoring_service

import (
	"strings"
	"sync"
	"time"

	andromodemError "github.com/basiooo/andromodem/internal/errors"
	"github.com/basiooo/andromodem/internal/model"
	"go.uber.org/zap"
)

type MonitoringTaskService struct {
	tasks  map[string]*model.MonitoringTask
	mutex  sync.RWMutex
	logger *zap.Logger
}

func NewMonitoringTaskService(logger *zap.Logger) IMonitoringTaskService {
	return &MonitoringTaskService{
		tasks:  make(map[string]*model.MonitoringTask),
		logger: logger,
	}
}

func (s *MonitoringTaskService) CreateTask(task *model.MonitoringTask) (*model.MonitoringTask, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.tasks[task.Serial]; exists {
		return nil, andromodemError.ErrorMonitoringTaskExists
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.IsActive = true

	s.tasks[task.Serial] = task

	s.logger.Info("[MonitoringTask] Created task",
		zap.String("serial", task.Serial),
		zap.String("host", task.Host),
		zap.String("method", string(task.Method)))

	return task, nil
}

func (s *MonitoringTaskService) UpdateTask(serial string, request *model.MonitoringTaskRequest) (*model.MonitoringTask, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.tasks[serial]
	if !exists {
		return nil, andromodemError.ErrorTaskNotFoundInConfig
	}

	if !s.ValidateTask(task) {
		return nil, andromodemError.ErrorInvalidMonitoringTask
	}

	tempTask := &model.MonitoringTask{
		Serial:            serial,
		Host:              request.Host,
		Method:            request.Method,
		MaxFailures:       request.MaxFailures,
		CheckingInterval:  request.CheckingInterval,
		AirplaneModeDelay: request.AirplaneModeDelay,
		CreatedAt:         task.CreatedAt,
		UpdatedAt:         time.Now(),
		IsActive:          task.IsActive,
	}

	if !s.ValidateTask(tempTask) {
		return nil, andromodemError.ErrorInvalidMonitoringTask
	}

	task.Host = request.Host
	task.Method = request.Method
	task.MaxFailures = request.MaxFailures
	task.CheckingInterval = request.CheckingInterval
	task.AirplaneModeDelay = request.AirplaneModeDelay
	task.UpdatedAt = time.Now()

	s.logger.Info("[MonitoringTask] Updated task",
		zap.String("serial", task.Serial),
		zap.String("host", task.Host),
		zap.String("method", string(task.Method)))

	return task, nil
}

func (s *MonitoringTaskService) UpdateTaskStatus(serial string, isActive bool) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.tasks[serial]
	if !exists {
		return andromodemError.ErrorTaskNotFoundInConfig
	}

	if !s.ValidateTask(task) {
		return andromodemError.ErrorInvalidMonitoringTask
	}

	task.IsActive = isActive
	task.UpdatedAt = time.Now()

	s.logger.Debug("[MonitoringTask] Updated task status",
		zap.String("serial", serial),
		zap.Bool("is_active", isActive))

	return nil
}

func (s *MonitoringTaskService) UpdateTaskField(serial string, updateFunc func(*model.MonitoringTask)) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, exists := s.tasks[serial]
	if !exists {
		return andromodemError.ErrorTaskNotFoundInConfig
	}

	if !s.ValidateTask(task) {
		return andromodemError.ErrorInvalidMonitoringTask
	}

	updateFunc(task)
	task.UpdatedAt = time.Now()

	if !s.ValidateTask(task) {
		return andromodemError.ErrorInvalidMonitoringTask
	}

	return nil
}

func (s *MonitoringTaskService) DeleteTask(serial string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.tasks[serial]; !exists {
		return andromodemError.ErrorTaskNotFoundInConfig
	}

	delete(s.tasks, serial)

	s.logger.Info("[MonitoringTask] Deleted task", zap.String("serial", serial))
	return nil
}

func (s *MonitoringTaskService) GetTask(serial string) (*model.MonitoringTask, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	task, exists := s.tasks[serial]
	if !exists {
		return nil, andromodemError.ErrorTaskNotFoundInConfig
	}

	return task, nil
}

func (s *MonitoringTaskService) GetAllTasks() ([]*model.MonitoringTask, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	tasks := make([]*model.MonitoringTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *MonitoringTaskService) ValidateTask(task *model.MonitoringTask) bool {
	if task == nil {
		return false
	}

	if strings.TrimSpace(task.Serial) == "" {
		return false
	}

	if strings.TrimSpace(task.Host) == "" {
		return false
	}

	if task.Method == "" {
		return false
	}

	if task.CheckingInterval <= 0 {
		return false
	}

	if task.AirplaneModeDelay < 0 {
		return false
	}

	validMethods := []model.MonitoringMethod{
		model.MethodPingByDevice,
		model.MethodHTTP,
		model.MethodHTTPS,
		model.MethodWS,
		model.MethodICMP,
	}

	for _, validMethod := range validMethods {
		if task.Method == validMethod {
			if task.MaxFailures < 1 || task.MaxFailures > 100 {
				return false
			}

			if !task.CreatedAt.IsZero() && task.CreatedAt.After(time.Now()) {
				return false
			}

			if !task.UpdatedAt.IsZero() && task.UpdatedAt.After(time.Now()) {
				return false
			}

			return true
		}
	}

	return false
}

func (s *MonitoringTaskService) TaskExists(serial string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists := s.tasks[serial]
	return exists
}

func (s *MonitoringTaskService) LoadTasks(tasks []*model.MonitoringTask) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, task := range tasks {
		s.tasks[task.Serial] = task
	}
}
