package service

import (
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

func NewDeviceService(adbClient *adb.Adb, adbCommand adbcommand.AdbCommand) DeviceService {
	return &DeviceServiceImpl{
		Adb:        adbClient,
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
		if err != nil {
			logrus.WithField("location", "DevicesService.GetDevices").Errorf("GetDevices(): failed get device state %v. error: %v", device.Model, err)
			state = goadb.StateInvalid
		}
		deviceState := adb.DeviceState(state)
		deviceResponse := model.Device{
			Serial: device.Serial,
			Model:  device.Model,
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
	deviceInfo := &model.DeviceInfo{}
	rootChan := make(chan parser.Root)
	batteryChan := make(chan parser.Battery)
	devicePropChan := make(chan parser.DeviceProp)
	go func() {
		defer close(rootChan)
		rootChan <- *d.GetRoot(*device)
	}()
	go func() {
		defer close(batteryChan)
		batteryChan <- *d.GetBattery(*device)
	}()
	go func() {
		defer close(devicePropChan)
		devicePropChan <- *d.GetDeviceProp(*device)
	}()
	deviceInfo.Root = <-rootChan
	deviceInfo.Battery = <-batteryChan
	deviceInfo.DeviceProp = <-devicePropChan
	return deviceInfo, nil
}
