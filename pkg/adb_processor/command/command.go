// Package command defines constants for various ADB commands used to interact with Android devices.
package command

type AdbCommand string

const (
	GetBatteryCommand                  AdbCommand = "dumpsys battery"                                                        //nolint:all, get battery info
	GetRootCommand                     AdbCommand = "su -v"                                                                  // get root version
	GetDevicePropCommand               AdbCommand = "getprop"                                                                // get all device property
	GetInboxCommand                    AdbCommand = "content query --uri content://sms/inbox --projection address,body,date" // get sms inbox message sort by date *Request root permission on android 9 and below
	GetAirplaneModeStatusLegacyCommand AdbCommand = "settings get global airplane_mode_on"                                   // get airplane mode status in Android 8 and below
	GetAirplaneModeStatusNewCommand    AdbCommand = "cmd connectivity airplane-mode"                                         // get airplane mode status in Android 9 and newer
	// Used in android 9 and newer
	EnableAirplaneModeNewCommand  AdbCommand = "cmd connectivity airplane-mode enable"  // enable airplane mode in Android 9 and newer
	DisableAirplaneModeNewCommand AdbCommand = "cmd connectivity airplane-mode disable" // disable airplane mode in Android 9 and newer
	// Used in android 8 and below
	EnableAirplaneModeLegacyCommand  AdbCommand = "settings put global airplane_mode_on 1" // enable airplane mode in Android 8 and below
	DisableAirplaneModeLegacyCommand AdbCommand = "settings put global airplane_mode_on 0" // disable airplane mode in Android 8 and below
	// Broadcast need root permission
	BroadcastAirplaneModeLegacyCommand AdbCommand = "am broadcast -a android.intent.action.AIRPLANE_MODE" // broadcast to save airplane mode in Android 8 and below

	GetMobileDataStateCommand AdbCommand = `dumpsys telephony.registry | grep -o 'mDataConnectionState=[^ ]*' | while read line; do 
		echo "${line#*=}" 
	done` // get mobile data state
	GetSimOperatorNameCommand     AdbCommand = "getprop gsm.sim.operator.alpha" // get sim card operator name from device property
	GetSimNetworkTypeCommand      AdbCommand = "getprop gsm.network.type"       // get sim network type from device property (usage for android 8 and below)
	GetDeviceModelCommand         AdbCommand = "getprop ro.product.model"
	GetDeviceProductCommand       AdbCommand = "getprop ro.product.name"
	GetAndroidVersionCommand      AdbCommand = "getprop ro.build.version.release"                                                                                 // get android version from device property
	GetSignalStrengthCommand      AdbCommand = "dumpsys telephony.registry | grep \"mSignalStrength=SignalStrength\" | sed 's/mSignalStrength=SignalStrength://'" // get signal strength
	GetApnCommand                 AdbCommand = "content query --uri content://telephony/carriers/preferapn"                                                       // get selected APN *not work on some devices, may need root access
	GetIpRouterCommand            AdbCommand = "ip route"                                                                                                         // get mobile data IP
	EnableMobileDataCommand       AdbCommand = "svc data enable"                                                                                                  // Enable mobile data (usage for android 8 and newer)
	DisableMobileDataCommand      AdbCommand = "svc data disable"                                                                                                 // Disable mobile data  (usage for android 8 and newer)
	GetMobileDataStatusCommand    AdbCommand = "settings get global mobile_data"                                                                                  // Get mobile data status
	RebootCommand                 AdbCommand = "reboot"                                                                                                           // reboot device
	RebootRecoveryCommand         AdbCommand = "reboot recovery"                                                                                                  // reboot device to recovery mode
	RebootBootloaderCommand       AdbCommand = "reboot bootloader"                                                                                                // reboot device to fastboot mode
	PowerOffCommand               AdbCommand = "reboot -p"                                                                                                        // power off device
	GetDeviceUptimeCommand        AdbCommand = "cat /proc/uptime"                                                                                                 // Get device uptime
	GetDeviceRootAccessCommand    AdbCommand = "su -c 'echo 1'"                                                                                                   // Get Root access
	GetDeviceMemoryCommand        AdbCommand = "cat /proc/meminfo"                                                                                                // Get Device Memory Info
	GetDeviceStorageCommand       AdbCommand = "dumpsys diskstats"                                                                                                // Get Internal Storage Info
	GetDeviceProcessNewCommand    AdbCommand = "ps -eo pid,user,%cpu,%mem,cmd,time+"                                                                              // Get Device Process (usage for android 8 and newer)
	GetDeviceProcessLegacyCommand AdbCommand = "ps"                                                                                                               // Get Device Process (usage for android 7 and below)
	GetBusyboxCheckCommand        AdbCommand = "which busybox"                                                                                                    // Check busybox installed or not
	GetKernelVersionCommand       AdbCommand = "uname -a"                                                                                                         // Get kernel version
)
