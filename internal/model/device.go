package model

import (
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
)

type Device struct {
	Serial   string `json:"serial"`
	Model    string `json:"model"`
	Product  string `json:"product"`
	State    string `json:"state,omitempty"`
	OldState string `json:"old_state,omitempty"`
	NewState string `json:"new_state,omitempty"`
}

type Other struct {
	Uptime            string `json:"uptime"`
	UptimeSecond      int    `json:"uptime_second"`
	BussyboxInstalled bool   `json:"busybox_installed"`
	KernelVersion     string `json:"kernel_version"`
}

type DeviceInfo struct {
	parser.DeviceProp `json:"prop"`
	parser.Root       `json:"root"`
	parser.Battery    `json:"battery"`
	Other             `json:"other"`
	parser.Memory     `json:"memory"`
	parser.Storage    `json:"storage"`
}

type FeatureAvailability struct {
	Key       string `json:"key"`
	Feature   string `json:"feature"`
	Available bool   `json:"available"`
	Message   string `json:"message"`
}

type FeatureAvailabilities struct {
	FeatureAvailabilities []FeatureAvailability `json:"feature_availabilities"`
}

type DeviceRootInfo struct {
	RootMethod  string
	Rooted      bool
	ShellAccess bool
}

type DeviceSpec struct {
	DeviceRootInfo
	AndroidVersion uint8
}

type DevicePowerAction struct {
	Action string `json:"action" validate:"required,oneof=power_off reboot reboot_recovery reboot_bootloader"`
}
