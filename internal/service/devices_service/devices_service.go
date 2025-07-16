package devices_service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	andromodemError "github.com/basiooo/andromodem/internal/errors"
	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/service/common_service"
	"github.com/basiooo/andromodem/pkg/adb_processor/command"
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/basiooo/andromodem/pkg/adb_processor/processor"
	"github.com/basiooo/andromodem/pkg/adb_processor/utils"
	"github.com/basiooo/andromodem/pkg/logger"
	adb "github.com/basiooo/goadb"
	"go.uber.org/zap"
)

type PowerAction string

const (
	PowerOff         PowerAction = "power_off"
	Reboot           PowerAction = "reboot"
	RebootRecovery   PowerAction = "reboot_recovery"
	RebootBootloader PowerAction = "reboot_bootloader"
)

type DevicesService struct {
	Adb          *adb.Adb
	AdbProcessor processor.IProcessor
	Logger       *zap.Logger
	Ctx          context.Context
}

func NewDevicesService(adb *adb.Adb, adbProcessor processor.IProcessor, logger *zap.Logger, ctx context.Context) IDevicesService {
	return &DevicesService{
		Adb:          adb,
		AdbProcessor: adbProcessor,
		Logger:       logger,
		Ctx:          ctx,
	}
}

func (d *DevicesService) DevicesListener(requestCtx context.Context, callback func(*model.Device) error) error {
	watcher := d.Adb.NewDeviceWatcher()
	eventChan := watcher.C()
	defer watcher.Shutdown()

	for {
		select {
		case <-requestCtx.Done():
			d.Logger.Info("Client closed the connection")
			return nil
		case <-d.Ctx.Done():
			d.Logger.Info("Server shutdown. Stopping device watcher")
			return d.Ctx.Err()
		case event, ok := <-eventChan:
			if !ok {
				d.Logger.Info("Watcher channel closed")
				return andromodemError.ErrorDevicesWatcherChannelClosed
			}
			if event.NewState == adb.StateOnline || event.NewState == adb.StateDisconnected || event.NewState == adb.StateAuthorizing {
				deviceModel := &model.Device{
					Serial:   event.Serial,
					OldState: event.OldState.String(),
					State:    event.NewState.String(),
					NewState: event.NewState.String(),
				}
				d.Logger.Debug("Device event", zap.String("serial", deviceModel.Serial), zap.String("old_state", deviceModel.OldState), zap.String("new_state", deviceModel.NewState))
				if event.NewState == adb.StateOnline {
					device, err := d.Adb.GetDeviceBySerial(event.Serial)
					if err == nil {
						if modelRequest, err := d.AdbProcessor.Run(device, command.GetDeviceModelCommand, false); err == nil {
							deviceModel.Model = utils.GetResultFromRaw(modelRequest)
						}
						if productRequest, err := d.AdbProcessor.Run(device, command.GetDeviceProductCommand, false); err == nil {
							deviceModel.Product = utils.GetResultFromRaw(productRequest)
						}
					}
				}
				err := callback(deviceModel)
				if err != nil {
					d.Logger.Error("error calling callback", zap.String("serial", deviceModel.Serial), zap.Error(err))
				}
			}
		}
	}
}

func (d *DevicesService) getDeviceMemory(device *adb.Device) (*parser.Memory, error) {
	defer logger.LogDuration(d.Logger, "GetDeviceInfo: get memory")()
	deviceMemory, err := d.AdbProcessor.Run(device, command.GetDeviceMemoryCommand, false)
	if err != nil {
		d.Logger.Error("error parsing device memory", zap.Error(err))
		return nil, err
	}
	if memory, ok := deviceMemory.(*parser.Memory); ok {
		return memory, nil
	}
	return nil, fmt.Errorf("error parsing device memory")
}

func (d *DevicesService) getDeviceStorage(device *adb.Device) (*parser.Storage, error) {
	defer logger.LogDuration(d.Logger, "GetDeviceInfo: get storage")()
	deviceStorage, err := d.AdbProcessor.Run(device, command.GetDeviceStorageCommand, false)
	if err != nil {
		d.Logger.Error("error parsing device storage", zap.Error(err))
		return nil, err
	}
	if storage, ok := deviceStorage.(*parser.Storage); ok {
		return storage, nil
	}
	return nil, fmt.Errorf("error parsing device storage")
}

func (d *DevicesService) getDeviceBattery(device *adb.Device) (*parser.Battery, error) {
	defer logger.LogDuration(d.Logger, "GetDeviceInfo: get battery")()
	deviceBattery, err := d.AdbProcessor.Run(device, command.GetBatteryCommand, false)
	if err != nil {
		d.Logger.Error("error parsing device battery", zap.Error(err))
		return nil, err
	}
	if battery, ok := deviceBattery.(*parser.Battery); ok {
		return battery, nil
	}
	return nil, fmt.Errorf("error parsing device battery")
}

func (d *DevicesService) getDeviceProps(device *adb.Device) (*parser.DeviceProp, error) {
	defer logger.LogDuration(d.Logger, "GetDeviceInfo: get props")()
	deviceProps, err := d.AdbProcessor.Run(device, command.GetDevicePropCommand, false)
	if err != nil {
		d.Logger.Error("error parsing device props", zap.Error(err))
		return nil, err
	}
	if props, ok := deviceProps.(*parser.DeviceProp); ok {
		return props, nil
	}
	return nil, fmt.Errorf("error parsing device props")
}

func (d *DevicesService) getDeviceRootInfo(device *adb.Device) (*parser.Root, error) {
	defer logger.LogDuration(d.Logger, "GetDeviceInfo: get root info")()
	root, err := d.AdbProcessor.RunWithRoot(device, command.GetRootCommand)
	if err != nil {
		d.Logger.Error("error parsing root", zap.Error(err))
		return nil, err
	}
	if root, ok := root.(*parser.Root); ok {
		return root, nil
	}
	return nil, fmt.Errorf("error parsing root")
}

func (d *DevicesService) getDeviceUptime(device *adb.Device) (*parser.DeviceUptime, error) {
	defer logger.LogDuration(d.Logger, "GetDeviceInfo: get uptime")()
	uptime, err := d.AdbProcessor.Run(device, command.GetDeviceUptimeCommand, false)
	if err != nil {
		d.Logger.Error("error parsing uptime", zap.Error(err))
		return nil, err
	}
	if uptime, ok := uptime.(*parser.DeviceUptime); ok {
		return uptime, nil
	}
	return nil, fmt.Errorf("error parsing uptime")
}

func (d *DevicesService) GetDeviceInfo(serial string) (*model.DeviceInfo, error) {
	defer logger.LogDuration(d.Logger, "GetDeviceInfo")()
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil || device == nil {
		d.Logger.Error("error getting device by serial", zap.String("serial", serial), zap.Error(err))
		return nil, andromodemError.ErrorDeviceNotFound
	}
	deviceInfo := &model.DeviceInfo{}
	var wg sync.WaitGroup
	wg.Add(7)
	go func() {
		defer wg.Done()
		memory, err := d.getDeviceMemory(device)
		if err != nil {
			d.Logger.Error("error parsing device memory", zap.String("serial", serial), zap.Error(err))
			return
		}
		deviceInfo.Memory = *memory
	}()
	go func() {
		defer wg.Done()
		storage, err := d.getDeviceStorage(device)
		if err != nil {
			d.Logger.Error("error parsing device storage", zap.String("serial", serial), zap.Error(err))
			return
		}
		deviceInfo.Storage = *storage
	}()
	go func() {
		defer wg.Done()
		props, err := d.getDeviceProps(device)
		if err != nil {
			d.Logger.Error("error parsing device props", zap.String("serial", serial), zap.Error(err))
			return
		}
		deviceInfo.DeviceProp = *props
	}()
	go func() {
		defer wg.Done()
		battery, err := d.getDeviceBattery(device)
		if err != nil {
			d.Logger.Error("error parsing device battery", zap.String("serial", serial), zap.Error(err))
			return
		}
		deviceInfo.Battery = *battery
	}()
	go func() {
		defer wg.Done()
		root, err := d.getDeviceRootInfo(device)
		if err != nil {
			d.Logger.Error("error parsing root", zap.String("serial", serial), zap.Error(err))
			return
		}
		deviceInfo.Root = *root
		if root.IsRooted {
			if shellRootAccess, err := d.AdbProcessor.RunWithRoot(device, command.GetDeviceRootAccessCommand); err == nil {
				deviceInfo.SuperUserAllowShellAccess = strings.Contains(utils.GetResultFromRaw(shellRootAccess.(*parser.RawParser)), "1")
			}
		}
	}()
	go func() {
		defer wg.Done()
		uptime, err := d.getDeviceUptime(device)
		if err != nil {
			d.Logger.Error("error parsing device uptime", zap.String("serial", serial), zap.Error(err))
			return
		}
		deviceInfo.Uptime = uptime.Uptime
		deviceInfo.UptimeSecond = uptime.UptimeSecond
	}()
	go func() {
		defer logger.LogDuration(d.Logger, "GetDeviceInfo: get kernel version")()
		defer wg.Done()
		_kernelVersion, err := d.AdbProcessor.Run(device, command.GetKernelVersionCommand, false)
		if err != nil {
			d.Logger.Error("error parsing device kernel version", zap.String("serial", serial), zap.Error(err))
		} else {
			deviceInfo.KernelVersion = utils.GetResultFromRaw(_kernelVersion)
		}
	}()
	wg.Wait()
	return deviceInfo, err
}

func (d *DevicesService) makeFeature(name, key string, available bool, availableMessage, unavailableMessage string) model.FeatureAvailability {
	message := availableMessage
	if !available {
		message = unavailableMessage
	}
	return model.FeatureAvailability{
		Feature:   name,
		Key:       key,
		Available: available,
		Message:   message,
	}
}

func (d *DevicesService) GetDeviceFeatureAvailabilities(serial string) (*model.FeatureAvailabilities, error) {
	defer logger.LogDuration(d.Logger, "GetDeviceFeatureAvailabilities")()
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil {
		d.Logger.Error("error getting device by serial", zap.String("serial", serial), zap.Error(err))
		return nil, andromodemError.ErrorDeviceNotFound
	}

	deviceSpec := &model.DeviceSpec{}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		rootInfo, err := common_service.GetDeviceRootAndAccessInfo(device, d.AdbProcessor, false)
		if err != nil {
			d.Logger.Error("error getting device root info", zap.String("serial", serial), zap.Error(err))
			return
		}
		deviceSpec.Rooted = rootInfo.Rooted
		deviceSpec.ShellAccess = rootInfo.ShellAccess
	}()
	go func() {
		defer logger.LogDuration(d.Logger, "GetDeviceFeatureAvailabilities: get android version")()
		defer wg.Done()
		if androidVersion, err := common_service.GetAndroidVersion(device, d.AdbProcessor, true); err == nil {
			deviceSpec.AndroidVersion = androidVersion
		}
	}()
	wg.Wait()
	features := []model.FeatureAvailability{
		d.makeFeature(
			"Can Read APN",
			"can_read_apn",
			deviceSpec.AndroidVersion >= command.MinimumAndroidShowApn,
			"Available",
			fmt.Sprintf("In Android %d or below, need root access", command.MinimumAndroidShowApn-1),
		),
		d.makeFeature(
			"Can Change Airplane Mode Status",
			"can_change_airplane_mode_status",
			deviceSpec.AndroidVersion >= command.MinimumAndroidToggleAirplaneMode,
			"Available",
			fmt.Sprintf("In Android %d or below, use broadcast command to change airplane mode status and need root access", command.MinimumAndroidToggleAirplaneMode-1),
		),
		d.makeFeature(
			"Can Change Sim Data Status",
			"can_change_sim_data_status",
			deviceSpec.AndroidVersion >= command.MinimumAndroidToggleMobileData,
			"Available (requires further verification on the actual device)", fmt.Sprintf("Only available in Android %d or below", command.MinimumAndroidToggleMobileData-1),
		),
		d.makeFeature(
			"Can Read Inbox",
			"can_read_inbox",
			deviceSpec.AndroidVersion >= command.MinimumAndroidShowMessages,
			"Available",
			fmt.Sprintf("In Android %d or below, need root access", command.MinimumAndroidShowMessages-1),
		),
		d.makeFeature(
			"Can Read Sim Signal Strength",
			"can_read_sim_signal_strength",
			deviceSpec.AndroidVersion >= command.MinimumAndroidGetSignalStrength,
			"Available",
			fmt.Sprintf("Only available in Android %d or above", command.MinimumAndroidGetSignalStrength),
		),
	}

	if deviceSpec.ShellAccess {
		for i := range features {
			switch features[i].Key {
			case "can_read_apn", "can_change_airplane_mode_status", "can_read_inbox":
				features[i].Available = true
				features[i].Message = ""
			}
		}
	}

	FeatureAvailabilities := &model.FeatureAvailabilities{
		FeatureAvailabilities: features,
	}
	return FeatureAvailabilities, nil
}

func (d *DevicesService) DevicePower(serial string, powerAction PowerAction) error {
	defer logger.LogDuration(d.Logger, "DevicePower")()
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil {
		d.Logger.Error("error getting device by serial", zap.String("serial", serial), zap.Error(err))
		return andromodemError.ErrorDeviceNotFound
	}
	var powerCommand command.AdbCommand
	switch powerAction {
	case PowerOff:
		powerCommand = command.PowerOffCommand
	case Reboot:
		powerCommand = command.RebootCommand
	case RebootRecovery:
		powerCommand = command.RebootRecoveryCommand
	case RebootBootloader:
		powerCommand = command.RebootBootloaderCommand
	}
	output, err := d.AdbProcessor.Run(device, powerCommand, false)
	if strings.TrimSpace(utils.GetResultFromRaw(output)) != "" {
		d.Logger.Error("error execute power action", zap.String("serial", serial), zap.String("action", string(powerAction)), zap.Error(err))
		return andromodemError.ErrorDevicePowerAction
	}
	return nil
}
