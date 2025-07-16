package monitoring_service

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/basiooo/andromodem/internal/model"
	"go.uber.org/zap"
)

type MonitoringConfigService struct {
	configFile   string
	logger       *zap.Logger
	taskService  IMonitoringTaskService
}

func NewMonitoringConfigService(configFile string, logger *zap.Logger, taskService IMonitoringTaskService) IMonitoringConfigService {
	return &MonitoringConfigService{
		configFile:  configFile,
		logger:      logger,
		taskService: taskService,
	}
}

func (s *MonitoringConfigService) LoadTasksFromFile() ([]*model.MonitoringTask, error) {
	if _, err := os.Stat(s.configFile); os.IsNotExist(err) {
		return []*model.MonitoringTask{}, nil
	}

	data, err := os.ReadFile(s.configFile)
	if err != nil {
		s.logger.Error("[MonitoringConfig] Failed to read config file, removing it", zap.Error(err))
		if removeErr := os.Remove(s.configFile); removeErr != nil {
			s.logger.Error("[MonitoringConfig] Failed to remove corrupted config file", zap.Error(removeErr))
		}
		return []*model.MonitoringTask{}, nil
	}

	if len(strings.TrimSpace(string(data))) == 0 {
		s.logger.Warn("[MonitoringConfig] Config file is empty, removing it")
		if removeErr := os.Remove(s.configFile); removeErr != nil {
			s.logger.Error("[MonitoringConfig] Failed to remove empty config file", zap.Error(removeErr))
		}
		return []*model.MonitoringTask{}, nil
	}

	var tasks []*model.MonitoringTask
	if err := json.Unmarshal(data, &tasks); err != nil {
		s.logger.Error("[MonitoringConfig] Config file contains invalid JSON, removing it", zap.Error(err))

		corruptedFile := s.configFile + ".corrupted." + time.Now().Format("20060102-150405")
		if copyErr := os.WriteFile(corruptedFile, data, 0644); copyErr != nil {
			s.logger.Warn("[MonitoringConfig] Failed to backup corrupted file", zap.Error(copyErr))
		}

		if removeErr := os.Remove(s.configFile); removeErr != nil {
			s.logger.Error("[MonitoringConfig] Failed to remove corrupted config file", zap.Error(removeErr))
		}
		return []*model.MonitoringTask{}, nil
	}

	validTasks := make([]*model.MonitoringTask, 0)
	invalidCount := 0

	for _, task := range tasks {
		if s.taskService.ValidateTask(task) {
			validTasks = append(validTasks, task)
		} else {
			invalidCount++
			serial := "unknown"
			host := "unknown"
			if task != nil {
				serial = task.Serial
				host = task.Host
			}
			s.logger.Warn("[MonitoringConfig] Found invalid task, skipping",
				zap.String("serial", serial),
				zap.String("host", host))
		}
	}

	if len(validTasks) == 0 {
		s.logger.Warn("[MonitoringConfig] No valid tasks found in config, removing file",
			zap.Int("total_tasks", len(tasks)),
			zap.Int("invalid_tasks", invalidCount))

		if removeErr := os.Remove(s.configFile); removeErr != nil {
			s.logger.Error("[MonitoringConfig] Failed to remove invalid config file", zap.Error(removeErr))
		}
		return []*model.MonitoringTask{}, nil
	}

	if invalidCount > 0 && float64(invalidCount)/float64(len(tasks)) > 0.5 {
		s.logger.Error("[MonitoringConfig] More than 50% of tasks are invalid, removing config file",
			zap.Int("total_tasks", len(tasks)),
			zap.Int("invalid_tasks", invalidCount),
			zap.Int("valid_tasks", len(validTasks)))

		if removeErr := os.Remove(s.configFile); removeErr != nil {
			s.logger.Error("[MonitoringConfig] Failed to remove mostly invalid config file", zap.Error(removeErr))
		}
		return []*model.MonitoringTask{}, nil
	}

	if invalidCount > 0 {
		s.logger.Info("[MonitoringConfig] Removed invalid tasks, saving cleaned config",
			zap.Int("original", len(tasks)),
			zap.Int("valid", len(validTasks)),
			zap.Int("removed", invalidCount))

		if saveErr := s.SaveTasksToFile(validTasks); saveErr != nil {
			s.logger.Error("[MonitoringConfig] Failed to save cleaned config", zap.Error(saveErr))
		}
	}

	return validTasks, nil
}

func (s *MonitoringConfigService) SaveTasksToFile(tasks []*model.MonitoringTask) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.configFile, data, 0644)
}

func (s *MonitoringConfigService) ValidateAndCleanConfig() error {
	if _, err := os.Stat(s.configFile); os.IsNotExist(err) {
		s.logger.Info("[MonitoringConfig] Config file does not exist, nothing to validate")
		return nil
	}

	data, err := os.ReadFile(s.configFile)
	if err != nil {
		s.logger.Error("[MonitoringConfig] Cannot read config file, removing it", zap.Error(err))
		return os.Remove(s.configFile)
	}

	if len(strings.TrimSpace(string(data))) == 0 {
		s.logger.Warn("[MonitoringConfig] Config file is empty, removing it")
		return os.Remove(s.configFile)
	}

	var tasks []*model.MonitoringTask
	if err := json.Unmarshal(data, &tasks); err != nil {
		s.logger.Error("[MonitoringConfig] Config file contains invalid JSON, removing it", zap.Error(err))
		return os.Remove(s.configFile)
	}

	validCount := 0
	for _, task := range tasks {
		if s.taskService.ValidateTask(task) {
			validCount++
		}
	}

	if validCount == 0 || float64(validCount)/float64(len(tasks)) < 0.5 {
		s.logger.Error("[MonitoringConfig] Config file contains mostly invalid data, removing it",
			zap.Int("total", len(tasks)),
			zap.Int("valid", validCount))
		return os.Remove(s.configFile)
	}

	s.logger.Info("[MonitoringConfig] Config file validation passed",
		zap.Int("total_tasks", len(tasks)),
		zap.Int("valid_tasks", validCount))

	return nil
}

func (s *MonitoringConfigService) SetConfigFile(configFile string) {
	s.configFile = configFile
}

func (s *MonitoringConfigService) GetConfigFile() string {
	return s.configFile
}