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
