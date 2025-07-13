package devices_service

import (
	"context"

	"github.com/basiooo/andromodem/internal/model"
)

type IDevicesService interface {
	DevicesListener(context.Context, func(*model.Device) error) error
	GetDeviceInfo(string) (*model.DeviceInfo, error)
	GetDeviceFeatureAvailabilities(string) (*model.FeatureAvailabilities, error)
	DevicePower(string, PowerAction) error
}
