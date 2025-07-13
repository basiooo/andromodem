package errors

import _errors "errors"

var (
	ErrorDeviceNotFound              = _errors.New("device not found")
	ErrorInvalidPowerAction          = _errors.New("invalid action")
	ErrorDevicesWatcherChannelClosed = _errors.New("devices watcher channel closed")
	ErrorDevicePowerAction           = _errors.New("error execute power action")
	// Network service errors
	ErrorCheckingAirplaneModeStatus = _errors.New("error checking airplane mode status")
	ErrorAirplaneModeActive         = _errors.New("airplane mode is currently active, unable to perform action, please disable airplane mode first")
	ErrorCheckingMobileDataState    = _errors.New("error checking mobile data state")
	ErrorTimeoutChangeMobileData    = _errors.New("timeout: cannot change mobile data state")
	ErrorTimeoutChangeAirplaneMode  = _errors.New("timeout: cannot change airplane mode state")
	// Monitoring service errors
	ErrorMonitoringTaskNotFound     = _errors.New("monitoring task not found")
	ErrorMonitoringTaskExists       = _errors.New("monitoring task already exists")
	ErrorNoRunningMonitoringTask    = _errors.New("no running monitoring task found")
	ErrorAutoStartTasksTimeout      = _errors.New("auto start tasks timeout")
	ErrorInvalidMonitoringTask      = _errors.New("invalid monitoring task configuration")
	ErrorTaskNotFoundInConfig       = _errors.New("task not found in configuration file")
	ErrorMonitoringTaskAlreadyRunning = _errors.New("monitoring task is already running")
)
