package adbcommand

import (
	goadb "github.com/abccyz/goadb"
	"github.com/sirupsen/logrus"
)

type AdbCommand struct {
}

func NewAdbCommand() *AdbCommand {
	return &AdbCommand{}
}

func (a *AdbCommand) GetBattery(device goadb.Device) (string, error) {
	batteryLevel, err := device.RunCommand("dumpsys battery")
	if err != nil {
		logrus.WithField("location", "AdbDeviceCommand.GetBattery").Errorf("GetBatteryLevel(): failed get device batery level: %v", err)
		return "", err
	}
	return batteryLevel, nil
}

func (a *AdbCommand) GetRoot(device goadb.Device) (string, error) {
	root, err := device.RunCommand(`su -v`)
	if err != nil {
		logrus.WithField("location", "AdbDeviceCommand.GetRoot").Errorf("GetRoot(): failed get device root status: %v", err)
		return "", err
	}
	return root, nil
}

func (a *AdbCommand) GetDeviceProp(device goadb.Device) (string, error) {
	deviceProps, err := device.RunCommand("getprop")
	if err != nil {
		logrus.WithField("location", "AdbDeviceCommand.GetDeviceProps").Errorf("GetDeviceProps(): failed get device props: %v", err)
		return "", err
	}
	return deviceProps, nil
}

func (a *AdbCommand) GetSmsInbox(device goadb.Device) (string, error) {
	inbox, err := device.RunCommand("content query --uri content://sms/inbox --projection address,body,date")
	if err != nil {
		logrus.WithField("location", "AdbDeviceCommand.GetSmsInbox").Errorf("GetSmsInbox(): failed get sms inbox: %v", err)
		return "", err
	}
	return inbox, nil
}

func (a *AdbCommand) GetSmsInboxRoot(device goadb.Device) (string, error) {
	inbox, err := device.RunCommand("su -c 'content query --uri content://sms/inbox --projection address,body,date'")
	if err != nil {
		logrus.WithField("location", "AdbDeviceCommand.GetSmsInboxRoot").Errorf("GetSmsInboxRoot(): failed get sms inbox with root method: %v", err)
		return "", err
	}
	return inbox, nil
}

func (a *AdbCommand) GetAirplaneModeStatus(device goadb.Device) (string, error) {
	airplaneModeStatus, err := device.RunCommand("cmd connectivity airplane-mode")
	if err != nil {
		logrus.WithField("location", "AdbDeviceCommand.GetAirplaneModeStatus").Errorf("GetAirplaneModeStatus(): failed get airplane mode status: %v", err)
		return "", err
	}
	return airplaneModeStatus, nil
}

func (a *AdbCommand) EnableAirplaneMode(device goadb.Device) (string, error) {
	enableAirplaneMode, err := device.RunCommand("cmd connectivity airplane-mode enable")
	if err != nil {
		logrus.WithField("location", "AdbDeviceCommand.EnableAirplaneMode").Errorf("EnableAirplaneMode(): failed enable airplane mode: %v", err)
		return "", err
	}
	return enableAirplaneMode, nil
}

func (a *AdbCommand) DisableAirplaneMode(device goadb.Device) (string, error) {
	disableAirplaneMode, err := device.RunCommand("cmd connectivity airplane-mode disable")
	if err != nil {
		logrus.WithField("location", "AdbDeviceCommand.DisableAirplaneMode").Errorf("DisableAirplaneMode(): failed enable airplane mode: %v", err)
		return "", err
	}
	return disableAirplaneMode, nil
}

func (a *AdbCommand) GetMobileDataState(device goadb.Device) (string, error) {
	mobileDataState, err := device.RunCommand("dumpsys telephony.registry | grep 'mDataConnectionState=' | cut -d'=' -f2 | grep -v '^$'")
	if err != nil {
		logrus.WithField("location", "AdbCommand.GetMobileDataState").Errorf("GetMobileDataState(): failed get device mobile data state: %v", err)
		return "", err
	}
	return mobileDataState, nil
}

func (a *AdbCommand) GetCarriersName(device goadb.Device) (string, error) {
	carriersName, err := device.RunCommand("getprop gsm.sim.operator.alpha")
	if err != nil {
		logrus.WithField("location", "AdbCommand.GetCarriersName").Errorf("GetCarriersName(): failed get device carriers name: %v", err)
		return "", err
	}
	return carriersName, nil
}

func (a *AdbCommand) GetSignalStrength(device goadb.Device) (string, error) {
	signalStrength, err := device.RunCommand(`dumpsys telephony.registry | grep "mSignalStrength=SignalStrength" | sed 's/mSignalStrength=SignalStrength://'`)
	if err != nil {
		logrus.WithField("location", "AdbCommand.GetSignalStrength").Errorf("GetSignalStrength(): failed get device signal strength: %v", err)
		return "", err
	}
	return signalStrength, nil
}

func (a *AdbCommand) GetApn(device goadb.Device) (string, error) {
	apn, err := device.RunCommand("content query --uri content://telephony/carriers/preferapn")
	if err != nil {
		logrus.WithField("location", "AdbCommand.GetApn").Errorf("GetApn(): failed get device Apn: %v", err)
		return "", err
	}
	return apn, nil
}

func (a *AdbCommand) GetMobileDataIp(device goadb.Device) (string, error) {
	mobileDataIp, err := device.RunCommand(`ip route | grep 'rmnet.*src' | sed -E 's/.*src ([^ ]+).*/\1/'`)
	if err != nil {
		logrus.WithField("location", "AdbCommand.GetMobileDataIp").Errorf("GetMobileDataIp(): failed get mobile data ip: %v", err)
		return "", err
	}
	return mobileDataIp, nil
}

func (a *AdbCommand) GetNetInterface(device goadb.Device, interfaceName string) (string, error) {
	interfaceData, err := device.RunCommand("ifconfig " + interfaceName)
	if err != nil {
		logrus.WithField("location", "AdbCommand.GetNetInterface").Errorf("GetNetInterface(): failed get network interface: %v", err)
		return "", err
	}
	return interfaceData, nil
}

func (a *AdbCommand) EnableMobileData(device goadb.Device) (string, error) {
	mobileData, err := device.RunCommand(`svc data enable`)
	if err != nil {
		logrus.WithField("location", "AdbCommand.EnableMobileData").Errorf("EnableMobileData(): failed enable mobile data: %v", err)
		return "", err
	}
	return mobileData, nil
}

func (a *AdbCommand) DisableMobileData(device goadb.Device) (string, error) {
	mobileData, err := device.RunCommand(`svc data disable`)
	if err != nil {
		logrus.WithField("location", "AdbCommand.DisableMobileData").Errorf("DisableMobileData(): failed enable mobile data: %v", err)
		return "", err
	}
	return mobileData, nil
}
