package handler

import (
	"context"
	adb "github.com/abccyz/goadb"
	adbcommands "github.com/basiooo/andromodem/adb_command"
	"github.com/basiooo/andromodem/helper"
	"github.com/basiooo/andromodem/model"
	"github.com/basiooo/andromodem/parser"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func fetchDeviceNetworkInfo(d *adb.Device) model.DeviceNetwork {
	deviceNetwork := model.DeviceNetwork{}
	airplaneModeStatusChan := make(chan bool)
	go func() {
		rawAirplaneModeStatus, _ := adbcommands.GetAirplaneModeStatus(d)
		airplaneModeStatusChan <- parser.ParseAirplaneModeStatus(rawAirplaneModeStatus)
		close(airplaneModeStatusChan)
	}()
	wifiChan := make(chan parser.Wifi)
	go func() {
		rawWifi, _ := adbcommands.GetWifi(d)
		wifiChan <- parser.ParseWifi(rawWifi)
		close(wifiChan)
	}()
	carriersChan := make(chan []parser.Carrier)
	go func() {
		rawCarrierData := parser.PrepareRawCarrier(d)
		carriersChan <- parser.ParseCarriers(rawCarrierData)
		close(carriersChan)
	}()
	deviceNetwork.Carriers = <-carriersChan
	deviceNetwork.AirplaneMode = <-airplaneModeStatusChan
	deviceNetwork.Wifi = <-wifiChan
	return deviceNetwork
}

func GetNetwork(writter http.ResponseWriter, request *http.Request) {
	adbClient, err := helper.GetADBClient(request)
	helper.PanicIfError(err)
	serial := chi.URLParam(request, "serial")
	device, err := helper.GetDeviceBySerial(adbClient, serial)
	if err != nil {
		logrus.WithField("function", "GetNetwork").Errorf("failed get device %v. error: %v", serial, err)
		response := model.WebResponseError{
			Error: "Device Not Found",
		}
		helper.WriteToResponseBody(writter, response, http.StatusNotFound)
		return
	}
	deviceNetwork := fetchDeviceNetworkInfo(device)
	helper.WriteToResponseBody(writter, deviceNetwork, http.StatusOK)
}

func deviceHasMobileDataEnable(device *adb.Device) bool {
	rawConnectionState, _ := adbcommands.GetMobileDataConnectionState(device)
	connectionState := parser.ParseMobileDataConnectionState(rawConnectionState)
	hasConnectionEnable := false
	for _, state := range connectionState {
		if state == parser.DataConnected || state == parser.DataConnecting {
			hasConnectionEnable = true
			break
		}
	}
	return hasConnectionEnable
}
func MobileDataToggle(writter http.ResponseWriter, request *http.Request) {
	adbClient, err := helper.GetADBClient(request)
	helper.PanicIfError(err)
	serial := chi.URLParam(request, "serial")
	device, err := helper.GetDeviceBySerial(adbClient, serial)
	if err != nil {
		logrus.WithField("function", "GetNetwork").Errorf("failed get device %v. error: %v", serial, err)
		response := model.WebResponseError{
			Error: "Device Not Found",
		}
		helper.WriteToResponseBody(writter, response, http.StatusNotFound)
		return
	}
	hasConnectionEnable := deviceHasMobileDataEnable(device)
	if hasConnectionEnable {
		err = adbcommands.DisableMobileData(device)
	} else {
		err = adbcommands.EnableMobileData(device)
	}
	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		if deviceHasMobileDataEnable(device) != hasConnectionEnable {
			break
		}
		select {
		case <-ctx.Done():
			break
		default:
		}
		time.Sleep(1 * time.Second)
	}
	response := model.ToggleMobileDataResponse{
		Success: true,
	}
	helper.PanicIfError(err)
	response.Enabled = !hasConnectionEnable
	helper.WriteToResponseBody(writter, response, http.StatusOK)
}
