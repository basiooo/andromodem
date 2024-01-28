package parser

type DeviceProps struct {
	Model          string `json:"model"`
	Brand          string `json:"brand"`
	Name           string `json:"name"`
	AndroidVersion string `json:"android_version"`
}

func ParseDeviceProps(rawDeviceProps string) DeviceProps {
	deviceProps := DeviceProps{}
	deviceProps.Model = GetPropValueByName(rawDeviceProps, "ro.product.model")
	deviceProps.Brand = GetPropValueByName(rawDeviceProps, "ro.product.brand")
	deviceProps.Name = GetPropValueByName(rawDeviceProps, "ro.product.name")
	deviceProps.AndroidVersion = GetPropValueByName(rawDeviceProps, "ro.build.version.release")
	return deviceProps
}
