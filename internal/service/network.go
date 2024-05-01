package service

import (
	"context"
	"errors"
	"strings"
	"sync"
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
	ToggleAirplaneMode(string) (*parser.AirplaneModeStatus, error)
	GetNetworkInfo(string) (*model.NetworkInfo, error)
	ToggleMobileData(string) (*model.ToggleMobileDataResponse, error)
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

func (d *NetworkServiceImpl) getMobileDataIp(device goadb.Device) *parser.IpAddress {
	rawMobileDataIp, _ := d.AdbCommand.GetMobileDataIp(device)
	return parser.NewIpAddress(rawMobileDataIp)
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

func (d *NetworkServiceImpl) ToggleAirplaneMode(serial string) (*parser.AirplaneModeStatus, error) {
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil {
		return nil, util.ErrDeviceNotFound
	}
	isEnabled := d.getAirplaneModeStatus(*device).Enabled
	result := &parser.AirplaneModeStatus{
		Enabled: isEnabled,
	}
	if isEnabled {
		_, err = d.AdbCommand.DisableAirplaneMode(*device)
	} else {
		_, err = d.AdbCommand.EnableAirplaneMode(*device)
	}
	if err == nil {
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
	outerLoop:
		for {
			if d.getAirplaneModeStatus(*device).Enabled != isEnabled {
				isEnabled = !isEnabled
				break outerLoop
			}
			select {
			case <-ctx.Done():
				err = errors.New("cannot change airplane mode state")
				break outerLoop
			default:
			}
			time.Sleep(1 * time.Second)
		}
	}
	result.Enabled = isEnabled
	return result, err
}

func (d *NetworkServiceImpl) getCarriers(device *goadb.Device) []parser.Carrier {
	var rawCarriers parser.RawCarriers
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		rawConnectionsState, _ := d.AdbCommand.GetMobileDataState(*device)
		rawCarriers.RawConnectionsState = rawConnectionsState
	}()
	go func() {
		defer wg.Done()
		rawCarriersName, _ := d.AdbCommand.GetCarriersName(*device)
		rawCarriers.RawCarriersName = rawCarriersName
	}()
	go func() {
		defer wg.Done()
		rawSignalsStrength, _ := d.AdbCommand.GetSignalStrength(*device)
		rawCarriers.RawSignalsStrength = rawSignalsStrength
	}()
	wg.Wait()
	carriersInfo := parser.NewCarriers(rawCarriers)
	return *carriersInfo
}

func (d *NetworkServiceImpl) GetNetworkInfo(serial string) (*model.NetworkInfo, error) {
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil {
		return nil, util.ErrDeviceNotFound
	}
	var networkInfo model.NetworkInfo
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		networkInfo.Carriers = d.getCarriers(device)
	}()
	go func() {
		defer wg.Done()
		networkInfo.AirplaneMode = d.getAirplaneModeStatus(*device).Enabled
	}()
	go func() {
		defer wg.Done()
		networkInfo.Apn = d.getApn(*device)
	}()
	go func() {
		defer wg.Done()
		ipInfo := d.getMobileDataIp(*device)
		networkInfo.Ip = ipInfo.Ip
	}()

	wg.Wait()

	return &networkInfo, nil

}

func (d *NetworkServiceImpl) deviceHasMobileDataEnable(device *goadb.Device) bool {
	rawMobileDataState, _ := d.AdbCommand.GetMobileDataState(*device)
	for _, data := range strings.Split(strings.TrimSpace(rawMobileDataState), "\n") {
		mobileDataState := parser.NewMobileDataConnectionState(data)
		state := mobileDataState.ConnectionState
		if state == parser.DataConnected || state == parser.DataConnecting {
			return true
		}
	}
	return false
}
func (d *NetworkServiceImpl) ToggleMobileData(serial string) (*model.ToggleMobileDataResponse, error) {
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil {
		return nil, util.ErrDeviceNotFound
	}
	isEnabled := d.deviceHasMobileDataEnable(device)
	airplaneModeStatusEnabled := d.getAirplaneModeStatus(*device).Enabled
	result := &model.ToggleMobileDataResponse{
		Enabled: isEnabled,
	}
	if airplaneModeStatusEnabled {
		err = errors.New("airplane mode is currently active, unable to perform action. Please disable airplane mode first.")
		return result, err
	}
	if isEnabled {
		_, err = d.AdbCommand.DisableMobileData(*device)
	} else {
		_, err = d.AdbCommand.EnableMobileData(*device)
	}

	if err == nil {
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
	outerLoop:
		for {
			if d.deviceHasMobileDataEnable(device) != isEnabled {
				isEnabled = !isEnabled
				break outerLoop
			}
			select {
			case <-ctx.Done():
				err = errors.New("cannot change mobile data state")
				break outerLoop
			default:
			}
			time.Sleep(1 * time.Second)
		}
	}

	result.Enabled = isEnabled
	return result, err
}
