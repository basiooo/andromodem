package network_service

import "github.com/basiooo/andromodem/internal/model"

type INetworkService interface {
	GetNetworkInfo(string) (*model.Network, error)
	ToggleMobileData(string) (*bool, error)
	ToggleAirplaneMode(string) (*bool, error)
}
