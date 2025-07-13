package parser

import (
	"strings"

	"github.com/basiooo/andromodem/pkg/adb_processor/errors"
)

type AirplaneModeState struct {
	Enabled bool `json:"enabled"`
}

func NewAirplaneModeState() IParser {
	return &AirplaneModeState{}
}
func (a *AirplaneModeState) Parse(rawData string) error {
	rawData = strings.TrimSpace(strings.ToLower(rawData))
	err := a.checkAvailability(rawData)
	if err != nil {
		return err
	}
	a.Enabled = isTruthy(rawData)
	return nil
}

func (a *AirplaneModeState) checkAvailability(rawData string) error {
	if strings.Contains(rawData, "no shell command implementation") {
		return errors.ErrorMinimumAndroidVersionNotSupport
	}
	return nil
}
