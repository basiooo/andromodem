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

func New() *Adb {
	return &Adb{}
}

func (a *Adb) Start() error {
	client, err := goadb.New()
	if err != nil {
		logrus.Errorf("Failed create ADB client: %v", err)
		return errors.New("failed create ADB client")
	}

	client.StartServer()
	logrus.Info("Trying start a ADB client manualy")
	cmd := exec.Command("adb", "start-server")
	err = cmd.Run()
	if err != nil {
		logrus.Errorf("Failed start ADB client manualy: %v", err)
		return errors.New("failed start ADB client manualy")
	}
	logrus.Info("Successfully run ADB client")
	a.Client = client
	return nil
}

func (a *Adb) GetDeviceBySerial(serial string) (*goadb.Device, error) {
	device := a.Client.Device(goadb.DeviceWithSerial(serial))
	_, err := device.Serial()
	return device, err
}

func (a *Adb) AdbIsInstalled() bool {
	cmd := exec.Command("which", "adb")
	err := cmd.Run()
	return err == nil
}
