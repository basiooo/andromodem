package monitoring_service

import (
	"context"
	"fmt"
	"sync"
	"time"

	andromodemError "github.com/basiooo/andromodem/internal/errors"
	"github.com/basiooo/andromodem/internal/model"
	"go.uber.org/zap"
)

type MonitoringWorkerService struct {
	runningTasks   map[string]context.CancelFunc
	status         map[string]*model.MonitoringStatus
	mutex          sync.RWMutex
	logger         *zap.Logger
	ctx            context.Context
	taskService    IMonitoringTaskService
	pinggerService IMonitoringPinggerService
	logService     IMonitoringLogService
	actionService  IDeviceActionService
	configService  IMonitoringConfigService
}

func NewMonitoringWorkerService(
	ctx context.Context,
	logger *zap.Logger,
	taskService IMonitoringTaskService,
	pinggerService IMonitoringPinggerService,
	logService IMonitoringLogService,
	actionService IDeviceActionService,
	configService IMonitoringConfigService,
) IMonitoringWorkerService {
	return &MonitoringWorkerService{
		runningTasks:   make(map[string]context.CancelFunc),
		status:         make(map[string]*model.MonitoringStatus),
		logger:         logger,
		ctx:            ctx,
		taskService:    taskService,
		pinggerService: pinggerService,
		logService:     logService,
		actionService:  actionService,
		configService:  configService,
	}
}

func (s *MonitoringWorkerService) saveTasksToFile() error {
	tasks, err := s.taskService.GetAllTasks()
	if err != nil {
		return err
	}

	return s.configService.SaveTasksToFile(tasks)
}

func (s *MonitoringWorkerService) StartMonitoring(serial string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, err := s.taskService.GetTask(serial)
	if err != nil {
		return andromodemError.ErrorTaskNotFoundInConfig
	}

	if !s.taskService.ValidateTask(task) {
		return andromodemError.ErrorInvalidMonitoringTask
	}

	if _, isRunning := s.runningTasks[serial]; isRunning {
		s.logger.Warn("[MonitoringWorker] Task already running", zap.String("serial", serial))
		return andromodemError.ErrorMonitoringTaskAlreadyRunning
	}

	s.logService.WriteLog(serial, true, "Monitoring task started")
	task.IsActive = true
	task.UpdatedAt = time.Now()

	s.status[serial] = &model.MonitoringStatus{
		Serial:       serial,
		IsRunning:    true,
		FailureCount: 0,
		LastSuccess:  false,
	}

	ctx, cancel := context.WithCancel(s.ctx)
	s.runningTasks[serial] = cancel

	go s.monitoringWorker(ctx, task)

	s.logger.Info("[MonitoringWorker] Started monitoring task",
		zap.String("serial", serial),
		zap.String("host", task.Host),
		zap.String("method", string(task.Method)))

	return nil
}

func (s *MonitoringWorkerService) StopMonitoring(serial string, isUpdate bool) error {
	if isUpdate {
		s.logService.WriteLog(serial, true, "Monitoring task stopped for update configuration")
	} else {
		s.logService.WriteLog(serial, true, "Monitoring task stopped")
	}

	task, err := s.taskService.GetTask(serial)
	if err != nil {
		return andromodemError.ErrorTaskNotFoundInConfig
	}

	if !s.taskService.ValidateTask(task) {
		return andromodemError.ErrorInvalidMonitoringTask
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	cancel, exists := s.runningTasks[serial]
	if !exists {
		return andromodemError.ErrorNoRunningMonitoringTask
	}

	cancel()
	delete(s.runningTasks, serial)

	if !isUpdate {
		if err := s.taskService.UpdateTaskStatus(serial, false); err != nil {
			s.logger.Warn("[MonitoringWorker] Failed to update task status",
				zap.String("serial", serial), zap.Error(err))
		}
	}

	delete(s.status, serial)

	s.logger.Info("[MonitoringWorker] Stopped monitoring task", zap.String("serial", serial))
	return nil
}

func (s *MonitoringWorkerService) GetStatus(serial string) (*model.MonitoringStatus, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	status, exists := s.status[serial]
	if !exists {
		return nil, andromodemError.ErrorDeviceNotFound
	}

	return status, nil
}

func (s *MonitoringWorkerService) IsRunning(serial string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists := s.runningTasks[serial]
	return exists
}

func (s *MonitoringWorkerService) AutoStartTasks() error {
	tasks, err := s.taskService.GetAllTasks()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		s.logger.Info("[MonitoringWorker] No monitoring tasks found, skipping auto start")
		return nil
	}

	activeTaskCount := 0
	var wg sync.WaitGroup
	errorChan := make(chan error, len(tasks))

	for _, task := range tasks {
		if task.IsActive {
			activeTaskCount++
			wg.Add(1)
			go func(t *model.MonitoringTask) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				select {
				case <-ctx.Done():
					errorChan <- fmt.Errorf("timeout starting task for serial %s", t.Serial)
					return
				default:
					s.logService.WriteLog(t.Serial, true, "Monitoring task automatically started on application startup")
					if err := s.StartMonitoring(t.Serial); err != nil {
						s.logger.Error("[MonitoringWorker] Failed to auto start monitoring task",
							zap.String("serial", t.Serial), zap.Error(err))
						errorChan <- err
					} else {
						s.logger.Info("[MonitoringWorker] Auto started monitoring task",
							zap.String("serial", t.Serial),
							zap.String("host", t.Host),
							zap.String("method", string(t.Method)))
					}
				}
			}(task)
		}
	}

	if activeTaskCount == 0 {
		s.logger.Info("[MonitoringWorker] No active monitoring tasks found, skipping auto start")
		return nil
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
		close(errorChan)
	}()

	select {
	case <-done:
		var errors []error
		for err := range errorChan {
			errors = append(errors, err)
		}
		if len(errors) > 0 {
			s.logger.Warn("[MonitoringWorker] Some tasks failed to start", zap.Int("failed_count", len(errors)))
		}
	case <-time.After(60 * time.Second):
		s.logger.Error("[MonitoringWorker] AutoStartTasks timeout, some tasks may not have started")
		return andromodemError.ErrorAutoStartTasksTimeout
	}

	return nil
}

func (s *MonitoringWorkerService) monitoringWorker(ctx context.Context, task *model.MonitoringTask) {
	ticker := time.NewTicker(time.Duration(task.CheckingInterval) * time.Second)
	defer ticker.Stop()

	failureCount := 0
	failureRestartCount := 0
	maxFailureRestartCount := 5

	var deviceOfflineStartTime *time.Time
	maxOfflineDuration := 5 * time.Minute

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Monitoring worker stopped", zap.String("serial", task.Serial))
			return
		case <-ticker.C:
			if !s.actionService.IsDeviceOnline(task.Serial) {
				s.logger.Debug("Device not online, waiting...", zap.String("serial", task.Serial))
				s.logService.WriteLog(task.Serial, false, "Device is offline, waiting for device to come online")

				if deviceOfflineStartTime == nil {
					now := time.Now()
					deviceOfflineStartTime = &now
					s.logger.Info("Device went offline, starting offline timer",
						zap.String("serial", task.Serial),
						zap.Time("offline_start_time", *deviceOfflineStartTime))
				}

				if time.Since(*deviceOfflineStartTime) >= maxOfflineDuration {
					s.logger.Error("Device offline for too long, stopping monitoring task",
						zap.String("serial", task.Serial),
						zap.Duration("offline_duration", time.Since(*deviceOfflineStartTime)),
						zap.Duration("max_offline_duration", maxOfflineDuration))
					s.logService.WriteLog(task.Serial, false, fmt.Sprintf("Device offline %v, stopping monitoring task", time.Since(*deviceOfflineStartTime).Round(time.Second)))

					if err := s.StopMonitoring(task.Serial, false); err != nil {
						s.logger.Error("[Monitoring] Failed to stop monitoring task after device offline timeout", zap.Error(err))
					}
					if err := s.saveTasksToFile(); err != nil {
						s.logger.Error("[Monitoring] Failed to save tasks to file after stop", zap.Error(err))
					}
					return
				}

				continue
			}

			if deviceOfflineStartTime != nil {
				s.logger.Info("Device back online, resetting offline timer",
					zap.String("serial", task.Serial),
					zap.Duration("was_offline_for", time.Since(*deviceOfflineStartTime)))
				deviceOfflineStartTime = nil
			}

			pingCtx, pingCancel := context.WithTimeout(ctx, 30*time.Second)
			success := s.pinggerService.PerformPing(pingCtx, task)
			pingCancel()

			s.mutex.Lock()
			if status, exists := s.status[task.Serial]; exists {
				status.LastPingTime = time.Now()
				status.LastSuccess = success
				if success {
					status.FailureCount = 0
					failureCount = 0
				} else {
					status.FailureCount++
					failureCount++
				}
			}
			s.mutex.Unlock()

			if success {
				message := fmt.Sprintf("Ping to %s success using %s method", task.Host, task.Method)
				s.logService.WriteLog(task.Serial, true, message)
				s.logger.Debug(message,
					zap.String("serial", task.Serial),
					zap.String("host", task.Host))
			} else {
				message := fmt.Sprintf("Ping to %s failed using %s method. Retry %d/%d", task.Host, task.Method, failureCount, task.MaxFailures)
				s.logService.WriteLog(task.Serial, false, message)
				s.logger.Debug(message,
					zap.String("serial", task.Serial),
					zap.String("host", task.Host),
					zap.Int("failure_count", failureCount))

				if failureCount >= task.MaxFailures {
					s.logger.Info("Max failures reached, performing restart action",
						zap.String("serial", task.Serial),
						zap.Int("failure_count", failureCount),
						zap.Int("max_failures", task.MaxFailures))

					for failureRestartCount <= maxFailureRestartCount {
						if err := s.actionService.PerformRestartAction(task.Serial, task.AirplaneModeDelay); err != nil {

							s.logger.Error("Failed to perform restart action",
								zap.String("serial", task.Serial),
								zap.Error(err))
							s.logService.WriteLog(task.Serial, false, fmt.Sprintf("Failed to perform restart action: %v. Retry %d/%d", err, failureRestartCount, maxFailureRestartCount))
							if failureRestartCount >= maxFailureRestartCount {
								s.logger.Error("Max restart failure reached, stopping monitoring task",
									zap.String("serial", task.Serial),
									zap.Int("failure_restart_count", failureRestartCount),
									zap.Int("max_failure_restart_count", maxFailureRestartCount))
								s.logService.WriteLog(task.Serial, false, "Max restart failure reached, stopping monitoring task")
								if err := s.StopMonitoring(task.Serial, false); err != nil {
									s.logger.Error("[Monitoring] Failed to stop monitoring task after max restart failure", zap.Error(err))
								}
								if err := s.saveTasksToFile(); err != nil {
									s.logger.Error("[Monitoring] Failed to save tasks to file after stop", zap.Error(err))
								}
								return
							}
							failureRestartCount++
							time.Sleep(5 * time.Second)
						} else {
							failureRestartCount = 0
							break
						}
					}
					s.logService.WriteLog(task.Serial, true, "Restart action performed successfully")
					failureCount = 0
					s.mutex.Lock()
					if status, exists := s.status[task.Serial]; exists {
						status.FailureCount = 0
					}
					s.mutex.Unlock()
				}
			}
		}
	}
}

func (s *MonitoringWorkerService) Shutdown(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logger.Info("[MonitoringWorker] Starting graceful shutdown")

	for serial, cancel := range s.runningTasks {
		cancel()
		s.logger.Info("[MonitoringWorker] Stopping monitoring task during shutdown",
			zap.String("serial", serial))

		if err := s.taskService.UpdateTaskStatus(serial, false); err != nil {
			s.logger.Warn("[MonitoringWorker] Failed to update task status during shutdown",
				zap.String("serial", serial), zap.Error(err))
		}

		if status, exists := s.status[serial]; exists {
			status.IsRunning = false
		}
	}

	s.runningTasks = make(map[string]context.CancelFunc)

	s.logger.Info("[MonitoringWorker] Graceful shutdown completed")
	return nil
}
