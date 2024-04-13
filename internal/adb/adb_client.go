package adb

import (
	"errors"
	"os/exec"

	goadb "github.com/abccyz/goadb"
	"github.com/sirupsen/logrus"
)

type Adb struct {
	Client *goadb.Adb
}

func NewAdbClient() *Adb {
	return &Adb{}
}

func (a *Adb) Start() error {
	client, err := goadb.New()
	if err != nil {
		logrus.Errorf("Failed create ADB client: %v", err)
		return errors.New("Failed create ADB client")
	}

	err = client.StartServer()
	if err != nil {
		logrus.Info("Trying start a ADB server manualy")
		cmd := exec.Command("adb", "start-server")
		err = cmd.Run()
		if err != nil {
			logrus.Errorf("Failed start ADB server manualy: %v", err)
			return errors.New("Failed start ADB server manualy")
		}
		logrus.Info("Successfully run ADB server")
	}
	a.Client = client
	return nil
}

func (a *Adb) GetDeviceBySerial(serial string) (*goadb.Device, error) {
	device := a.Client.Device(goadb.DeviceWithSerial(serial))
	_, err := device.Serial()
	return device, err
}
