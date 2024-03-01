package service

import (
	"sync"

	goadb "github.com/abccyz/goadb"
	"github.com/basiooo/andromodem/internal/adb"
	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/util"
	adbcommand "github.com/basiooo/andromodem/pkg/adb/adb_command"
	"github.com/basiooo/andromodem/pkg/adb/parser"
	"github.com/sirupsen/logrus"
)

type DeviceService interface {
	GetDevices() ([]model.Device, error)
	GetDeviceInfo(string) (*model.DeviceInfo, error)
}

type DeviceServiceImpl struct {
	*adb.Adb
	AdbCommand adbcommand.AdbCommand
}

func NewDeviceService(adb *adb.Adb, adbCommand adbcommand.AdbCommand) DeviceService {
	return &DeviceServiceImpl{
		Adb:        adb,
		AdbCommand: adbCommand,
	}
}

func (d *DeviceServiceImpl) GetDevices() ([]model.Device, error) {
	deviceList, err := d.Adb.Client.ListDevices()
	if err != nil {
		return nil, err
	}
	devices := make([]model.Device, 0, len(deviceList))
	for _, device := range deviceList {
		deviceInfo := d.Adb.Client.Device(goadb.DeviceWithSerial(device.Serial))
		state, err := deviceInfo.State()
		device_props := d.GetDeviceProp(*deviceInfo)
		if err != nil {
			logrus.WithField("location", "DevicesService.GetDevices").Errorf("GetDevices(): failed get device state %v. error: %v", device.Model, err)
			state = goadb.StateInvalid
		}
		deviceState := adb.DeviceState(state)
		deviceResponse := model.Device{
			Serial: device.Serial,
			Model:  device_props.Model,
			State:  deviceState.String(),
		}
		devices = append(devices, deviceResponse)
	}
	return devices, nil
}

func (d *DeviceServiceImpl) GetRoot(device goadb.Device) *parser.Root {
	rawRoot, _ := d.AdbCommand.GetRoot(device)
	root := parser.NewRoot(rawRoot)
	return root
}

func (d *DeviceServiceImpl) GetBattery(device goadb.Device) *parser.Battery {
	rawBattery, _ := d.AdbCommand.GetBattery(device)
	battery := parser.NewBattery(rawBattery)
	return battery
}

func (d *DeviceServiceImpl) GetDeviceProp(device goadb.Device) *parser.DeviceProp {
	rawDeviceProp, _ := d.AdbCommand.GetDeviceProp(device)
	deviceProp := parser.NewDeviceProp(rawDeviceProp)
	return deviceProp
}

func (d *DeviceServiceImpl) GetDeviceInfo(serial string) (*model.DeviceInfo, error) {
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil {
		return nil, util.ErrDeviceNotFound
	}
	var deviceInfo model.DeviceInfo
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		deviceInfo.Root = *d.GetRoot(*device)
	}()
	go func() {
		defer wg.Done()
		deviceInfo.Battery = *d.GetBattery(*device)
	}()
	go func() {
		defer wg.Done()
		deviceInfo.DeviceProp = *d.GetDeviceProp(*device)
	}()
	wg.Wait()
	return &deviceInfo, nil
}
