package service

import (
	"context"
	"errors"
	"time"

	goadb "github.com/abccyz/goadb"
	"github.com/basiooo/andromodem/internal/adb"
	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/util"
	adbcommand "github.com/basiooo/andromodem/pkg/adb/adb_command"
	"github.com/basiooo/andromodem/pkg/adb/parser"
)

type NetworkService interface {
	GetAirplaneModeStatus(string) (*parser.AirplaneModeStatus, error)
	ToggleAirplaneMode(string) (*model.AirplaneModeResponse, error)
	GetNetworkInfo(string) (*model.NetworkInfo, error)
}

type NetworkServiceImpl struct {
	*adb.Adb
	AdbCommand adbcommand.AdbCommand
}

func NewNetworkService(adb *adb.Adb, adbCommand adbcommand.AdbCommand) NetworkService {
	return &NetworkServiceImpl{
		Adb:        adb,
		AdbCommand: adbCommand,
	}
}

func (d *NetworkServiceImpl) getAirplaneModeStatus(device goadb.Device) *parser.AirplaneModeStatus {
	rawAirplaneModeStatus, _ := d.AdbCommand.GetAirplaneModeStatus(device)
	airplaneModeStatus := parser.NewAirplaneModeStatus(rawAirplaneModeStatus)
	return airplaneModeStatus
}

func (d *NetworkServiceImpl) getApn(device goadb.Device) *parser.Apn {
	rawApn, _ := d.AdbCommand.GetApn(device)
	apn := parser.NewApn(rawApn)
	return apn
}

func (d *NetworkServiceImpl) EnableAirplaneMode(device goadb.Device) error {
	_, err := d.AdbCommand.EnableAirplaneMode(device)
	return err
}

func (d *NetworkServiceImpl) DisableAirplaneMode(device goadb.Device) error {
	_, err := d.AdbCommand.DisableAirplaneMode(device)
	return err
}

func (d *NetworkServiceImpl) GetAirplaneModeStatus(serial string) (*parser.AirplaneModeStatus, error) {
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil {
		return nil, util.ErrDeviceNotFound
	}
	airplaneModeStatusChain := make(chan *parser.AirplaneModeStatus)
	go func() {
		defer close(airplaneModeStatusChain)
		airplaneModeStatusChain <- d.getAirplaneModeStatus(*device)
	}()

	return <-airplaneModeStatusChain, nil
}

func (d *NetworkServiceImpl) ToggleAirplaneMode(serial string) (*model.AirplaneModeResponse, error) {
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil {
		return nil, util.ErrDeviceNotFound
	}
	isEnabled := d.getAirplaneModeStatus(*device).Enabled
	isSuccess := true
	if isEnabled {
		err = d.DisableAirplaneMode(*device)
	} else {
		err = d.EnableAirplaneMode(*device)
	}
	if err == nil {
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
	outerLoop:
		for {
			if d.getAirplaneModeStatus(*device).Enabled != isEnabled {
				isEnabled = !isEnabled
				isSuccess = true
				break outerLoop
			}
			select {
			case <-ctx.Done():
				err = errors.New("cannot change airplane mode state")
				isSuccess = false
				break outerLoop
			default:
			}
			time.Sleep(1 * time.Second)
		}
	}
	res := model.AirplaneModeResponse{}
	res.Enabled = isEnabled
	res.Success = isSuccess
	if err != nil {
		res.Error = err.Error()
	}
	return &res, err
}

func (d *NetworkServiceImpl) getCarriers(device *goadb.Device) []parser.Carrier {
	rawConnectionsState, _ := d.AdbCommand.GetMobileDataState(*device)
	rawCarriersName, _ := d.AdbCommand.GetCarriersName(*device)
	rawSignalsStrength, _ := d.AdbCommand.GetSignalStrength(*device)
	rawCarriers := parser.RawCarriers{
		RawConnectionsState: rawConnectionsState,
		RawCarriersName:     rawCarriersName,
		RawSignalsStrength:  rawSignalsStrength,
	}
	carriersInfo := parser.NewCarriers(rawCarriers)
	return *carriersInfo
}

func (d *NetworkServiceImpl) GetNetworkInfo(serial string) (*model.NetworkInfo, error) {
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil {
		return nil, util.ErrDeviceNotFound
	}
	networkInfo := model.NetworkInfo{}
	carriersChan := make(chan []parser.Carrier)
	apnChan := make(chan *parser.Apn)
	go func() {
		defer close(carriersChan)
		carriersChan <- d.getCarriers(device)
	}()
	go func() {
		defer close(apnChan)
		apnChan <- d.getApn(*device)
	}()
	networkInfo.Carriers = <-carriersChan
	networkInfo.Apn = *<-apnChan
	return &networkInfo, nil
}
