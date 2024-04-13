package model

type BaseToggleResponse struct {
	Success bool   `json:"success"`
	Enabled bool   `json:"enabled"`
	Error   string `json:"errors,omitempty"`
}

type AirplaneModeResponse struct {
	Success bool   `json:"success"`
	Enabled bool   `json:"enabled"`
	Error   string `json:"errors,omitempty"`
}
