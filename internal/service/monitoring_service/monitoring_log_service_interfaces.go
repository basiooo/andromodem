package monitoring_service

import (
	"context"

	"github.com/basiooo/andromodem/internal/model"
)

type IMonitoringLogService interface {
	WriteLog(string, bool, string)
	GetLogs(string, int) ([]*model.MonitoringLog, error)
	SetLogDir(string)
	GetLogDir() string
	LogListener(ctx context.Context, serial string, callback func(*model.MonitoringLog) error) error
	NotifyNewLog(serial string, log *model.MonitoringLog)
	ClearLogs(string) error
	ClearAllLogs() error
	Shutdown()
}
