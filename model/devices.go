package model

type Devices struct {
	Serial string `json:"serial"`
	Model  string `json:"model"`
	State  string `json:"state"`
}

type DevicesResponse struct {
	Devices []Devices `json:"devices"`
}
