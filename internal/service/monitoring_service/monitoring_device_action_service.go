package monitoring_service

import (
	"time"

	network_service "github.com/basiooo/andromodem/internal/service/network"
	adb "github.com/basiooo/goadb"
	"go.uber.org/zap"
)

type MonitoringDeviceActionService struct {
	adb            *adb.Adb
	networkService network_service.INetworkService
	logService     IMonitoringLogService
	logger         *zap.Logger
}

func NewMonitoringDeviceActionService(
	adb *adb.Adb,
	networkService network_service.INetworkService,
	logService IMonitoringLogService,
	logger *zap.Logger,
) IDeviceActionService {
	return &MonitoringDeviceActionService{
		adb:            adb,
		networkService: networkService,
		logService:     logService,
		logger:         logger,
	}
}

func (s *MonitoringDeviceActionService) IsDeviceOnline(serial string) bool {
	device, err := s.adb.GetDeviceBySerial(serial)
	if err != nil || device == nil {
		return false
	}
	state, _ := device.State()
	return state == adb.StateOnline
}

func (s *MonitoringDeviceActionService) PerformRestartAction(serial string, airplaneModeDelay int) error {
	s.logger.Info("Performing restart action: toggle mobile data", zap.String("serial", serial))
	first_state, err := s.networkService.ToggleAirplaneMode(serial)
	if err != nil {
		s.logger.Error("Failed to toggle airplane mode", zap.String("serial", serial), zap.Error(err))
		s.logService.WriteLog(serial, false, "Failed to enable airplane mode during restart action")
		return err
	}
	if first_state != nil && *first_state {
		s.logService.WriteLog(serial, true, "Success enable airplane mode")
	}

	if first_state == nil || *first_state {
		if airplaneModeDelay > 0 {
			time.Sleep(time.Duration(airplaneModeDelay) * time.Second)
		}
		second_state, err := s.networkService.ToggleAirplaneMode(serial)
		if err != nil {
			s.logger.Error("Failed to toggle airplane mode", zap.String("serial", serial), zap.Error(err))
			s.logService.WriteLog(serial, false, "Failed to disable airplane mode during restart action")
			return err
		}
		if second_state != nil && !*second_state {
			s.logService.WriteLog(serial, true, "Success disable airplane mode")
		}
	}
	return nil
}
