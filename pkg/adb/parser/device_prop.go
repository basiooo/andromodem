package parser

import (
	"fmt"
	"regexp"
)

type DeviceProp struct {
	Model          string `json:"model"`
	Brand          string `json:"brand"`
	Name           string `json:"name"`
	AndroidVersion string `json:"android_version"`
}

func NewDeviceProp(rawDeviceProp string) *DeviceProp {
	deviceProps := DeviceProp{}
	deviceProps.Model = GetPropValueByName(rawDeviceProp, "ro.product.model")
	deviceProps.Brand = GetPropValueByName(rawDeviceProp, "ro.product.brand")
	deviceProps.Name = GetPropValueByName(rawDeviceProp, "ro.product.name")
	deviceProps.AndroidVersion = GetPropValueByName(rawDeviceProp, "ro.build.version.release")
	return &deviceProps
}

func GetPropValueByName(prop, name string) string {
	regexPattern := fmt.Sprintf(`\[%s\]: \[(.*?)\]`, name)
	r := regexp.MustCompile(regexPattern)

	match := r.FindStringSubmatch(prop)

	if len(match) < 2 {
		return ""
	}

	return match[1]
}
