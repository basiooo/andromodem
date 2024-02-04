package handler

import (
	adbcommands "github.com/basiooo/andromodem/adb_command"
	"github.com/basiooo/andromodem/parser"
	"net/http"

	adb "github.com/abccyz/goadb"
	"github.com/basiooo/andromodem/helper"
	"github.com/basiooo/andromodem/model"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func fetchDeviceInfo(d *adb.Device) model.DeviceInfo {
	deviceInfo := model.DeviceInfo{}
	rootChan := make(chan parser.Root)
	go func() {
		rawRoot, _ := adbcommands.GetRoot(d)
		rootChan <- parser.ParseRoot(rawRoot)
		close(rootChan)
	}()
	devicePropsChan := make(chan parser.DeviceProps)
	go func() {
		rawDeviceProps, _ := adbcommands.GetDeviceProps(d)
		devicePropsChan <- parser.ParseDeviceProps(rawDeviceProps)
		close(devicePropsChan)
	}()
	batteryLevelChan := make(chan string)
	go func() {
		rawBatteryLevel, _ := adbcommands.GetBatteryLevel(d)
		batteryLevelChan <- parser.ParseBatteryLevel(rawBatteryLevel)
		close(batteryLevelChan)
	}()
	deviceInfo.Root = <-rootChan
	deviceInfo.DeviceProps = <-devicePropsChan
	deviceInfo.BatteryLevel = <-batteryLevelChan
	return deviceInfo
}

func GetDevice(writter http.ResponseWriter, request *http.Request) {
	adbClient, err := helper.GetADBClient(request)
	helper.PanicIfError(err)
	serial := chi.URLParam(request, "serial")
	device, err := helper.GetDeviceBySerial(adbClient, serial)
	if err != nil {
		logrus.WithField("function", "GetDevice").Errorf("failed get device %v. error: %v", serial, err)
		response := model.WebResponseError{
			Error: "Device Not Found",
		}
		helper.WriteToResponseBody(writter, response, http.StatusNotFound)
		return
	}
	deviceInfo := fetchDeviceInfo(device)
	helper.WriteToResponseBody(writter, deviceInfo, http.StatusOK)

}

func GetThermal(writter http.ResponseWriter, request *http.Request) {
	adbClient, err := helper.GetADBClient(request)
	helper.PanicIfError(err)
	serial := chi.URLParam(request, "serial")
	device, err := helper.GetDeviceBySerial(adbClient, serial)
	if err != nil {
		logrus.WithField("function", "GetDeviceThermal").Errorf("failed get device %v. error: %v", serial, err)
		response := model.WebResponseError{
			Error: "Device Not Found",
		}
		helper.WriteToResponseBody(writter, response, http.StatusNotFound)
		return
	}
	rawThermal, err := adbcommands.GetDeviceThermal(device)
	if err != nil {
		logrus.WithField("function", "GetDeviceThermal").Errorf("failed get thermal %v. error: %v", serial, err)
	}
	thermal := parser.ParseThermal(rawThermal)
	deviceThermal := thermal.Data
	helper.WriteToResponseBody(writter, deviceThermal, http.StatusOK)

}
