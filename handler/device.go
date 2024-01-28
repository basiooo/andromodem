package handler

import (
	adb "github.com/abccyz/goadb"
	adbcommands "github.com/basiooo/andromodem/adb_command"
	"github.com/basiooo/andromodem/helper"
	"github.com/basiooo/andromodem/model"
	"github.com/basiooo/andromodem/parser"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
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
