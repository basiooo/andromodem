package processor

import (
	"github.com/basiooo/andromodem/pkg/adb_processor/command"
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
)

var commandParser = map[command.AdbCommand]any{
	command.GetDevicePropCommand:               parser.NewDeviceProp,
	command.GetAndroidVersionCommand:           parser.NewRawParser,
	command.GetDeviceModelCommand:              parser.NewRawParser,
	command.GetDeviceProductCommand:            parser.NewRawParser,
	command.GetBatteryCommand:                  parser.NewBattery,
	command.GetRootCommand:                     parser.NewRoot,
	command.GetInboxCommand:                    parser.NewInbox,
	command.GetApnCommand:                      parser.NewApn,
	command.GetAirplaneModeStatusLegacyCommand: parser.NewAirplaneModeState,
	command.GetAirplaneModeStatusNewCommand:    parser.NewAirplaneModeState,
	command.EnableAirplaneModeNewCommand:       parser.NewRawParser,
	command.DisableAirplaneModeNewCommand:      parser.NewRawParser,
	command.EnableAirplaneModeLegacyCommand:    parser.NewRawParser,
	command.DisableAirplaneModeLegacyCommand:   parser.NewRawParser,
	command.BroadcastAirplaneModeLegacyCommand: parser.NewRawParser,
	command.EnableMobileDataCommand:            parser.NewRawParser,
	command.DisableMobileDataCommand:           parser.NewRawParser,
	command.GetMobileDataStatusCommand:         parser.NewRawParser,
	command.RebootCommand:                      parser.NewRawParser,
	command.RebootRecoveryCommand:              parser.NewRawParser,
	command.RebootBootloaderCommand:            parser.NewRawParser,
	command.PowerOffCommand:                    parser.NewRawParser,
	command.GetDeviceUptimeCommand:             parser.NewDeviceUptime,
	command.GetDeviceRootAccessCommand:         parser.NewRawParser,
	command.GetIpRouterCommand:                 parser.NewIPRoute,
	command.GetKernelVersionCommand:            parser.NewRawParser,
	//  keep for future
	//	command.GetSimOperatorNameCommand:       parser.NewDeviceSim,
	//	command.GetSignalStrengthCommand:        parser.NewSignalStrength,
	//	command.GetMobileDataStateCommand:       parser.NewMobileDataState,

	command.GetSimOperatorNameCommand: parser.NewRawParser,
	command.GetSignalStrengthCommand:  parser.NewRawParser,
	command.GetMobileDataStateCommand: parser.NewRawParser,

	command.GetSimNetworkTypeCommand: parser.NewRawParser,
	command.GetDeviceMemoryCommand:   parser.NewMemory,
	command.GetDeviceStorageCommand:  parser.NewStorage,
}
