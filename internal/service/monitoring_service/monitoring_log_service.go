package monitoring_service

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/basiooo/andromodem/internal/model"
	"go.uber.org/zap"
)

type MonitoringLogService struct {
	logDir        string
	logger        *zap.Logger
	logMutex      sync.RWMutex
	logBuffers    map[string][]string
	maxLines      int
	logListeners  map[string][]chan *model.MonitoringLog
	listenerMutex sync.RWMutex

	writeQueue map[string]*writeQueueItem
	writeMutex sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

type writeQueueItem struct {
	writeChan   chan struct{}
	writeCancel context.CancelFunc
}

func NewMonitoringLogService(logDir string, logger *zap.Logger) IMonitoringLogService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &MonitoringLogService{
		logDir:       logDir,
		logger:       logger,
		logBuffers:   make(map[string][]string),
		maxLines:     100,
		logListeners: make(map[string][]chan *model.MonitoringLog),
		writeQueue:   make(map[string]*writeQueueItem),
		ctx:          ctx,
		cancel:       cancel,
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		logger.Error("[MonitoringLog] Failed to create log directory", zap.Error(err))
	}

	return service
}

func (s *MonitoringLogService) WriteLog(serial string, success bool, message string) {
	log := &model.MonitoringLog{
		Serial:    serial,
		Success:   success,
		Message:   message,
		Timestamp: time.Now(),
	}

	logData, err := json.Marshal(log)
	if err != nil {
		s.logger.Error("[MonitoringLog] Failed to marshal log", zap.Error(err))
		return
	}

	logLine := string(logData)

	s.logMutex.Lock()
	if s.logBuffers[serial] == nil {
		s.logBuffers[serial] = make([]string, 0, s.maxLines)
	}

	s.logBuffers[serial] = append(s.logBuffers[serial], logLine)
	if len(s.logBuffers[serial]) > s.maxLines {
		s.logBuffers[serial] = s.logBuffers[serial][1:]
	}
	s.logMutex.Unlock()

	s.NotifyNewLog(serial, log)

	s.writeBufferToFile(serial)
}

func (s *MonitoringLogService) Shutdown() {
	s.cancel()
	s.wg.Wait()

	s.listenerMutex.Lock()
	for serial, listeners := range s.logListeners {
		for _, ch := range listeners {
			close(ch)
		}
		delete(s.logListeners, serial)
	}
	s.listenerMutex.Unlock()

	s.writeMutex.Lock()
	for serial, item := range s.writeQueue {
		item.writeCancel()
		close(item.writeChan)
		delete(s.writeQueue, serial)
	}
	s.writeMutex.Unlock()
}

func (s *MonitoringLogService) writeBufferToFile(serial string) {
	s.logMutex.RLock()
	buffer := make([]string, len(s.logBuffers[serial]))
	copy(buffer, s.logBuffers[serial])
	s.logMutex.RUnlock()

	if len(buffer) == 0 {
		return
	}

	logFile := filepath.Join(s.logDir, fmt.Sprintf("%s.log", serial))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		s.logger.Error("[MonitoringLog] Failed to open log file", zap.Error(err))
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.logger.Error("Failed to close log file", zap.Error(err))
		}
	}()

	for _, line := range buffer {
		if _, err := file.WriteString(line + "\n"); err != nil {
			s.logger.Error("[MonitoringLog] Failed to write to log file", zap.Error(err))
			return
		}
	}
}

func (s *MonitoringLogService) GetLogs(serial string, limit int) ([]*model.MonitoringLog, error) {
	logFile := filepath.Join(s.logDir, fmt.Sprintf("%s.log", serial))

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return []*model.MonitoringLog{}, nil
	}

	file, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.logger.Error("Failed to close log file", zap.Error(err))
		}
	}()

	var logs []*model.MonitoringLog
	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	start := 0
	if limit > 0 && len(lines) > limit {
		start = len(lines) - limit
	}

	for i := start; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		var log model.MonitoringLog
		if err := json.Unmarshal([]byte(line), &log); err != nil {
			s.logger.Error("Failed to parse log line", zap.Error(err), zap.String("line", line))
			continue
		}

		logs = append(logs, &log)
	}

	return logs, nil
}

func (s *MonitoringLogService) ClearLogs(serial string) error {
	s.logMutex.Lock()
	delete(s.logBuffers, serial)
	s.logMutex.Unlock()

	logFile := filepath.Join(s.logDir, fmt.Sprintf("%s.log", serial))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		s.logger.Error("[MonitoringLog] Failed to clear log file", zap.Error(err))
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.logger.Error("Failed to close log file", zap.Error(err))
		}
	}()

	s.logger.Info("[MonitoringLog] Cleared logs for serial", zap.String("serial", serial))
	return nil
}

func (s *MonitoringLogService) ClearAllLogs() error {
	s.logMutex.Lock()
	s.logBuffers = make(map[string][]string)
	s.logMutex.Unlock()

	files, err := filepath.Glob(filepath.Join(s.logDir, "*.log"))
	if err != nil {
		s.logger.Error("[MonitoringLog] Failed to list log files", zap.Error(err))
		return err
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			s.logger.Error("[MonitoringLog] Failed to remove log file", zap.Error(err), zap.String("file", file))
			return err
		}
	}

	s.logger.Info("[MonitoringLog] Cleared all logs")
	return nil
}

func (s *MonitoringLogService) SetLogDir(logDir string) {
	s.logDir = logDir
	if err := os.MkdirAll(logDir, 0755); err != nil {
		s.logger.Error("[MonitoringLog] Failed to create log directory", zap.Error(err))
	}
}

func (s *MonitoringLogService) GetLogDir() string {
	return s.logDir
}

func (s *MonitoringLogService) LogListener(ctx context.Context, serial string, callback func(*model.MonitoringLog) error) error {
	logChan := make(chan *model.MonitoringLog, 10)

	s.listenerMutex.Lock()
	if s.logListeners[serial] == nil {
		s.logListeners[serial] = make([]chan *model.MonitoringLog, 0)
	}
	s.logListeners[serial] = append(s.logListeners[serial], logChan)
	s.listenerMutex.Unlock()

	defer func() {
		s.listenerMutex.Lock()
		defer s.listenerMutex.Unlock()

		if listeners, exists := s.logListeners[serial]; exists {
			for i, ch := range listeners {
				if ch == logChan {
					s.logListeners[serial] = append(listeners[:i], listeners[i+1:]...)
					break
				}
			}
			if len(s.logListeners[serial]) == 0 {
				delete(s.logListeners, serial)
			}
		}
		close(logChan)
		s.logger.Debug("[MonitoringLog] Listener cleaned up", zap.String("serial", serial))
	}()

	existingLogs, err := s.GetLogs(serial, s.maxLines)
	if err != nil {
		s.logger.Error("[MonitoringLog] Failed to get existing logs", zap.Error(err))
	} else {
		for _, log := range existingLogs {
			if err := callback(log); err != nil {
				return err
			}
		}
	}

	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("[MonitoringLog] Context cancelled, stopping listener", zap.String("serial", serial))
			return ctx.Err()

		case log, ok := <-logChan:
			if !ok {
				return nil
			}
			if err := callback(log); err != nil {
				return err
			}
		}
	}
}

func (s *MonitoringLogService) NotifyNewLog(serial string, log *model.MonitoringLog) {
	s.listenerMutex.RLock()
	listeners, exists := s.logListeners[serial]
	if !exists {
		s.listenerMutex.RUnlock()
		return
	}

	listenersCopy := make([]chan *model.MonitoringLog, len(listeners))
	copy(listenersCopy, listeners)
	s.listenerMutex.RUnlock()

	for _, ch := range listenersCopy {
		select {
		case ch <- log:
		default:
			s.logger.Debug("[MonitoringLog] Listener channel full/closed, skipping", zap.String("serial", serial))
		}
	}
}
