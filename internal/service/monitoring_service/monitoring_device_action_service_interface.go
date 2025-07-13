package monitoring_service

type IDeviceActionService interface {
	PerformRestartAction(string, int) error
	IsDeviceOnline(string) bool
}
