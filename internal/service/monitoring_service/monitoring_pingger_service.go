package monitoring_service

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/basiooo/andromodem/internal/model"
	adb "github.com/basiooo/goadb"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type MonitoringPinggerService struct {
	adb        *adb.Adb
	httpClient *http.Client
	logger     *zap.Logger
}

func NewMonitoringPinggerService(adb *adb.Adb, logger *zap.Logger) IMonitoringPinggerService {
	return &MonitoringPinggerService{
		adb:    adb,
		logger: logger,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (s *MonitoringPinggerService) PerformPing(ctx context.Context, task *model.MonitoringTask) bool {
	switch task.Method {
	case model.MethodPingByDevice:
		return s.PingByDevice(ctx, task.Serial, task.Host)
	case model.MethodHTTP:
		return s.PingHTTP(ctx, task.Host, false)
	case model.MethodHTTPS:
		return s.PingHTTP(ctx, task.Host, true)
	case model.MethodWS:
		return s.PingWebSocket(ctx, task.Host)
	case model.MethodICMP:
		return s.PingICMP(ctx, task.Host)
	default:
		s.logger.Error("Unknown monitoring method", zap.String("method", string(task.Method)))
		return false
	}
}

func (s *MonitoringPinggerService) PingHTTP(ctx context.Context, host string, useHTTPS bool) bool {
	scheme := "http"
	if useHTTPS {
		scheme = "https"
	}

	url := fmt.Sprintf("%s://%s", scheme, host)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			if closeErr := resp.Body.Close(); closeErr != nil {
				s.logger.Error("Failed to close HTTP response body", zap.Error(closeErr))
			}
		}
	}()

	return true
}

func (s *MonitoringPinggerService) PingWebSocket(ctx context.Context, host string) bool {
	url := fmt.Sprintf("ws://%s", host)

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return false
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			s.logger.Error("Failed to close WebSocket connection", zap.Error(closeErr))
		}
	}()

	return true
}

func (s *MonitoringPinggerService) isIMCPSuccess(result string) bool {
	if len(result) == 0 {
		return false
	}

	resultLower := strings.ToLower(result)

	successPatterns := []string{
		"1 received",        // Linux
		" 0% packet loss",   // Linux/macOS
		" 0.0% packet loss", // macOS
		"received = 1",      // Windows
		"lost = 0",          // Windows
	}

	for _, pattern := range successPatterns {
		if strings.Contains(resultLower, pattern) {
			return true
		}
	}

	return false
}

func (s *MonitoringPinggerService) PingByDevice(ctx context.Context, serial, host string) bool {
	device, err := s.adb.GetDeviceBySerial(serial)
	if err != nil {
		s.logger.Error("Failed to get device", zap.String("serial", serial), zap.Error(err))
		return false
	}

	cmd := fmt.Sprintf("ping -c 1 -W 5 %s", host)
	result, err := device.RunCommand(cmd)
	if err != nil {
		s.logger.Error("Failed to run ping command", zap.Error(err))
		return false
	}

	return s.isIMCPSuccess(result)
}

func (s *MonitoringPinggerService) PingICMP(ctx context.Context, host string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "ping", "-n", "1", "-w", "5000", host)
	case "darwin":
		cmd = exec.CommandContext(ctx, "ping", "-c", "1", "-W", "5000", host)
	default:
		cmd = exec.CommandContext(ctx, "ping", "-c", "1", "-W", "5", host)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Debug("Failed to run ping command", zap.String("host", host), zap.Error(err))
	}

	return s.isIMCPSuccess(string(output))
}
