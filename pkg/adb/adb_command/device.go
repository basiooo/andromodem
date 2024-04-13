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
