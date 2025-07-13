package monitoring_service

import (
	"context"

	"github.com/basiooo/andromodem/internal/model"
)

type IMonitoringPinggerService interface {
	PerformPing(context.Context, *model.MonitoringTask) bool
	PingByDevice(context.Context, string, string) bool
	PingHTTP(context.Context, string, bool) bool
	PingWebSocket(context.Context, string) bool
	PingICMP(context.Context, string) bool
}
