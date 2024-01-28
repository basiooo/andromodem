package adbcommands

import (
	"log"

	adb "github.com/abccyz/goadb"
)

func GetWifi(d *adb.Device) (string, error) {
	wifi, err := d.RunCommand("dumpsys wifi")
	if err != nil {
		log.Printf("GetWifi(): %v", err)
		return "", err
	}
	return wifi, nil
}
func GetCarrierNames(d *adb.Device) (string, error) {
	carrierNames, err := d.RunCommand("getprop gsm.sim.operator.alpha")
	if err != nil {
		log.Printf("GetCarrierNames(): %v", err)
		return "", err
	}
	return carrierNames, nil
}
func GetAirplaneModeStatus(d *adb.Device) (string, error) {
	airplaneModeStatus, err := d.RunCommand("cmd connectivity airplane-mode")
	if err != nil {
		log.Printf("GetAirplaneModeStatus(): %v", err)
		return "", err
	}
	return airplaneModeStatus, nil
}

func GetMobileDataConnectionState(d *adb.Device) (string, error) {
	mobileDataConnectionState, err := d.RunCommand("dumpsys telephony.registry | grep 'mDataConnectionState=' | cut -d'=' -f2 | grep -v '^$'")
	if err != nil {
		log.Printf("GetMobileDataConnectionState(): %v", err)
		return "", err
	}
	return mobileDataConnectionState, nil
}

func GetMobileDataSignalStrength(d *adb.Device) (string, error) {
	mobileDataSignalStrength, err := d.RunCommand(`dumpsys telephony.registry | grep "mSignalStrength=SignalStrength" | sed 's/mSignalStrength=SignalStrength://'`)
	if err != nil {
		log.Printf("GetMobileDataSignalStrength(): %v", err)
		return "", err
	}
	return mobileDataSignalStrength, nil
}
func GetBatteryLevel(d *adb.Device) (string, error) {
	batteryLevel, err := d.RunCommand("dumpsys battery | grep level:  | sed 's/.*/&%/' | cut -d : -f2")
	if err != nil {
		log.Printf("GetBatteryLevel(): %v", err)
		return "", err
	}
	return batteryLevel, nil
}

func GetDeviceProps(d *adb.Device) (string, error) {
	deviceProps, err := d.RunCommand("getprop")
	if err != nil {
		log.Printf("GetDeviceProps(): %v", err)
		return "", err
	}
	return deviceProps, nil
}

func GetRoot(d *adb.Device) (string, error) {
	isRooted, err := d.RunCommand(`su -v`)
	if err != nil {
		log.Printf("GetRoot(): %v", err)
		return "", err
	}
	return isRooted, nil
}

func EnableMobileData(d *adb.Device) error {
	_, err := d.RunCommand(`svc data enable`)
	if err != nil {
		log.Printf("EnableMobileData(): %v", err)
	}
	return err
}

func DisableMobileData(d *adb.Device) error {
	_, err := d.RunCommand(`svc data disable`)
	if err != nil {
		log.Printf("DisableMobileData(): %v", err)
	}
	return err
}

func GetDeviceTemp(d *adb.Device) (string, error) {
	deviceTemp, err := d.RunCommand(`
	for VARIABLE in $(seq 0 100)
	do
		if [ -e /sys/devices/virtual/thermal/thermal_zone$VARIABLE/type ]; then
		cat /sys/devices/virtual/thermal/thermal_zone$VARIABLE/type
		cat /sys/devices/virtual/thermal/thermal_zone$VARIABLE/temp
		fi
	done
	`)
	if err != nil {
		log.Printf("GetDeviceTemp(): %v", err)
		return "", err
	}
	return deviceTemp, nil
}
