package helper

import (
	adb "github.com/abccyz/goadb"
)

func GetDeviceBySerial(adbClient *adb.Adb, serial string) (*adb.Device, error) {
	device := adbClient.Device(adb.DeviceWithSerial(serial))
	_, err := device.Serial()
	return device, err
}
