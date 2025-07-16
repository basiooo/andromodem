package network_service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	adbErrors "github.com/basiooo/andromodem/pkg/adb_processor/errors"

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

type NetworkService struct {
	Adb          *adb.Adb
	AdbProcessor processor.IProcessor
	Logger       *zap.Logger
	Ctx          context.Context
}

func NewNetworkService(adb *adb.Adb, adbProcessor processor.IProcessor, logger *zap.Logger, ctx context.Context) INetworkService {
	return &NetworkService{
		Adb:          adb,
		AdbProcessor: adbProcessor,
		Logger:       logger,
		Ctx:          ctx,
	}
}

func (n *NetworkService) getIpRoutes(device *adb.Device) (*parser.IPRoute, error) {
	defer logger.LogDuration(n.Logger, "getIpRoutes")()
	ipRoutes, err := n.AdbProcessor.Run(device, command.GetIpRouterCommand, true)
	if err != nil {
		n.Logger.Error("error parsing ip routes", zap.Error(err))
		return nil, err
	}
	if result, ok := ipRoutes.(*parser.IPRoute); ok {
		return result, nil
	}
	return nil, fmt.Errorf("error parsing ip routes")
}

func (n *NetworkService) getApn(device *adb.Device) (*parser.Apn, error) {
	defer logger.LogDuration(n.Logger, "getApn")()
	needRoot := false
	var apn parser.IParser
	var err error
	if androidVersion, err := common_service.GetAndroidVersion(device, n.AdbProcessor, true); err == nil {
		needRoot = androidVersion < command.MinimumAndroidShowApn
	}
	if needRoot {
		apn, err = n.AdbProcessor.RunWithRoot(device, command.GetApnCommand)
	} else {
		apn, err = n.AdbProcessor.Run(device, command.GetApnCommand, false)
	}
	if err != nil {
		n.Logger.Error("error parsing apn", zap.Error(err))
		return nil, err
	}
	if result, ok := apn.(*parser.Apn); ok {
		return result, nil
	}
	return nil, fmt.Errorf("error parsing apn")
}

func (n *NetworkService) getAirplaneModeStatus(device *adb.Device) (bool, error) {
	defer logger.LogDuration(n.Logger, "getAirplaneModeStatus")()
	cmd := command.GetAirplaneModeStatusNewCommand
	if androidVersion, err := common_service.GetAndroidVersion(device, n.AdbProcessor, true); err == nil {
		if androidVersion < command.MinimumAndroidToggleAirplaneMode {
			cmd = command.GetAirplaneModeStatusLegacyCommand
		}
	}
	airplaneModeStatus, err := n.AdbProcessor.Run(device, cmd, false)
	if err != nil {
		n.Logger.Error("error parsing airplane mode status", zap.Error(err))
		return false, err
	}
	if result, ok := airplaneModeStatus.(*parser.AirplaneModeState); ok {
		return result.Enabled, nil
	}
	return false, fmt.Errorf("error parsing airplane mode status")
}

func (n *NetworkService) getSimsRaw(device *adb.Device) (parser.RawDeviceSim, error) {
	defer logger.LogDuration(n.Logger, "getSimsRaw")()
	var rawDeviceSim parser.RawDeviceSim
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		if rawOperatorName, err := n.AdbProcessor.Run(device, command.GetSimOperatorNameCommand, false); err == nil {
			rawDeviceSim.RawCarriersName = utils.GetResultFromRaw(rawOperatorName)
		}
	}()
	go func() {
		defer wg.Done()
		if rawConnectionState, err := n.AdbProcessor.Run(device, command.GetMobileDataStateCommand, false); err == nil {
			rawDeviceSim.RawConnectionsState = utils.GetResultFromRaw(rawConnectionState)
		}
	}()
	go func() {
		defer wg.Done()
		cmd := command.GetSignalStrengthCommand
		if androidVersion, err := common_service.GetAndroidVersion(device, n.AdbProcessor, true); err == nil {
			if androidVersion < command.MinimumAndroidGetSignalStrength {
				// in android 9 or bellow cannot parse signal strength
				cmd = command.GetSimNetworkTypeCommand
			}
		}
		if rawSignalsStrength, err := n.AdbProcessor.Run(device, cmd, false); err == nil {
			rawDeviceSim.RawSignalsStrength = utils.GetResultFromRaw(rawSignalsStrength)
		}
	}()
	wg.Wait()
	return rawDeviceSim, nil
}

func (n *NetworkService) GetNetworkInfo(serial string) (*model.Network, error) {
	defer logger.LogDuration(n.Logger, "GetNetworkInfo")()
	var err error
	device, err := n.Adb.GetDeviceBySerial(serial)
	if err != nil || device == nil {
		n.Logger.Error("error getting device by serial", zap.String("serial", serial), zap.Error(err))
		return nil, andromodemError.ErrorDeviceNotFound
	}
	networkInfo := &model.Network{}
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		ipRoutes, err := n.getIpRoutes(device)
		if err != nil {
			n.Logger.Error("error getting ip routes", zap.String("serial", serial), zap.Error(err))
			return
		}
		if ipRoutes.NetworkIPs != nil {
			networkInfo.IpRoutes = ipRoutes.NetworkIPs
		} else {
			networkInfo.IpRoutes = []parser.NetworkIp{}
		}
	}()
	go func() {
		defer wg.Done()
		apn, err := n.getApn(device)
		if err != nil {
			n.Logger.Error("error getting apn", zap.String("serial", serial), zap.Error(err))
			return
		}
		networkInfo.APN = *apn
	}()
	go func() {
		defer wg.Done()
		airplaneMode, err := n.getAirplaneModeStatus(device)
		if err != nil {
			n.Logger.Error("error getting airplane mode status", zap.String("serial", serial), zap.Error(err))
			return
		}
		networkInfo.AirplaneMode = airplaneMode
	}()
	go func() {
		// TODO: refactor to use parser interface
		defer wg.Done()
		rawSims, err := n.getSimsRaw(device)
		if err != nil {
			n.Logger.Error("error getting sims", zap.String("serial", serial), zap.Error(err))
			return
		}
		deviceSims := parser.NewDeviceSim()
		data, err := json.Marshal(rawSims)
		if err != nil {
			n.Logger.Error("error marshalling sims", zap.String("serial", serial), zap.Error(err))
			return
		}
		err = deviceSims.Parse(string(data))
		if err != nil {
			n.Logger.Error("error parsing sims", zap.String("serial", serial), zap.Error(err))
			return
		}
		networkInfo.Sims = deviceSims.(*parser.DeviceSim).Sims
	}()
	wg.Wait()
	return networkInfo, nil
}
func (n *NetworkService) hasMobileDataEnabled(device *adb.Device) (bool, error) {
	defer logger.LogDuration(n.Logger, "hasMobileDataEnabled")()

	rawConnectionState, err := n.AdbProcessor.Run(device, command.GetMobileDataStateCommand, false)
	if err != nil {
		return false, err
	}

	stateParser := parser.NewMobileDataState()
	for _, state := range strings.Split(strings.TrimSpace(utils.GetResultFromRaw(rawConnectionState)), "\n") {
		if err := stateParser.Parse(state); err != nil {
			serial, err := device.Serial()
			if err != nil {
				n.Logger.Error("error getting device serial", zap.Error(err))
				return false, err
			}
			n.Logger.Error("error parsing mobile data state", zap.String("serial", serial), zap.Error(err))
			return false, err
		}
		if currentState := stateParser.(*parser.MobileDataState).State; currentState == parser.DataConnected || currentState == parser.DataConnecting {
			return true, nil
		}
	}

	return false, nil
}
func (n *NetworkService) ToggleMobileData(serial string) (*bool, error) {
	defer logger.LogDuration(n.Logger, "ToggleMobileData")()

	device, err := n.Adb.GetDeviceBySerial(serial)
	if err != nil || device == nil {
		n.Logger.Error("error getting device by serial", zap.String("serial", serial), zap.Error(err))
		return nil, andromodemError.ErrorDeviceNotFound
	}

	isAirplaneMode, err := n.getAirplaneModeStatus(device)
	if err != nil {
		n.Logger.Error("error checking airplane mode status", zap.String("serial", serial), zap.Error(err))
		return nil, andromodemError.ErrorCheckingAirplaneModeStatus
	}
	if isAirplaneMode {
		return nil, andromodemError.ErrorAirplaneModeActive
	}

	isDataEnabled, err := n.hasMobileDataEnabled(device)
	if err != nil {
		n.Logger.Error("error checking mobile data state", zap.String("serial", serial), zap.Error(err))
		return nil, andromodemError.ErrorCheckingMobileDataState
	}

	var cmd command.AdbCommand
	newState := "enable"
	if isDataEnabled {
		newState = "disable"
		cmd = command.DisableMobileDataCommand
	} else {
		cmd = command.EnableMobileDataCommand
	}

	if _, err := n.AdbProcessor.Run(device, cmd, false); err != nil {
		n.Logger.Error("error toggling mobile data", zap.String("serial", serial), zap.Error(err))
		return nil, fmt.Errorf("error %s mobile data: %w", newState, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, andromodemError.ErrorTimeoutChangeMobileData
		case <-ticker.C:
			currentStatus, err := n.hasMobileDataEnabled(device)
			if err != nil {
				n.Logger.Error("error checking mobile data state", zap.String("serial", serial), zap.Error(err))
				return nil, andromodemError.ErrorCheckingMobileDataState
			}
			if currentStatus != isDataEnabled {
				return &currentStatus, nil
			}
		}
	}
}

func (n *NetworkService) ToggleAirplaneMode(serial string) (*bool, error) {
	defer logger.LogDuration(n.Logger, "ToggleAirplaneMode")()

	device, err := n.Adb.GetDeviceBySerial(serial)
	if err != nil || device == nil {
		n.Logger.Error("error getting device by serial", zap.String("serial", serial), zap.Error(err))
		return nil, andromodemError.ErrorDeviceNotFound
	}

	useLegacyCommand := false
	if androidVersion, err := common_service.GetAndroidVersion(device, n.AdbProcessor, true); err == nil {
		if androidVersion < command.MinimumAndroidToggleAirplaneMode {
			rootInfo, err := common_service.GetDeviceRootAndAccessInfo(device, n.AdbProcessor, false)
			if err != nil {
				return nil, err
			}
			if !rootInfo.ShellAccess {
				return nil, adbErrors.ErrorNeedRoot
			}
			useLegacyCommand = true
		}
	}

	isEnabled, err := n.getAirplaneModeStatus(device)
	if err != nil {
		n.Logger.Error("error checking airplane mode status", zap.String("serial", serial), zap.Error(err))
		return nil, andromodemError.ErrorCheckingAirplaneModeStatus
	}

	var cmd command.AdbCommand
	newState := "enable"
	if isEnabled {
		newState = "disable"
		if useLegacyCommand {
			cmd = command.DisableAirplaneModeLegacyCommand
		} else {
			cmd = command.DisableAirplaneModeNewCommand
		}
	} else {
		if useLegacyCommand {
			cmd = command.EnableAirplaneModeLegacyCommand
		} else {
			cmd = command.EnableAirplaneModeNewCommand
		}
	}

	if _, err := n.AdbProcessor.Run(device, cmd, false); err != nil {
		n.Logger.Error("error toggling airplane mode", zap.String("serial", serial), zap.Error(err))
		return nil, fmt.Errorf("error %s airplane mode: %w", newState, err)
	}

	if useLegacyCommand {
		if _, err := n.AdbProcessor.RunWithRoot(device, command.BroadcastAirplaneModeLegacyCommand); err != nil {
			n.Logger.Error("error broadcasting airplane mode", zap.String("serial", serial), zap.Error(err))
			return nil, fmt.Errorf("error broadcasting airplane mode: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, andromodemError.ErrorTimeoutChangeAirplaneMode
		case <-ticker.C:
			currentStatus, err := n.getAirplaneModeStatus(device)
			if err != nil {
				n.Logger.Error("error checking airplane mode status", zap.String("serial", serial), zap.Error(err))
				return nil, andromodemError.ErrorCheckingAirplaneModeStatus
			}
			if currentStatus != isEnabled {
				return &currentStatus, nil
			}
		}
	}
}
