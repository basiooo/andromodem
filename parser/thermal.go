package parser

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type Thermal struct {
	Data map[string]int
}

func toCelcius(data string) int {
	value, err := strconv.Atoi(data)
	if err != nil {
		logrus.WithField("function", "toCelcius").Error(err)
		return 0
	}
	if value > 10000 {
		value = value / 1000
	} else if value > 1000 {
		value = value / 100
	} else if value > 100 {
		value = value / 10
	}
	return value
}
func ParseThermal(rawThermal string) Thermal {
	thermal := Thermal{}
	lines := strings.Split(rawThermal, "\n")
	thermalMap := make(map[string]int)
	var key string
	for i, line := range lines {
		if (i % 2) == 0 {
			key = line
		} else {
			thermalMap[key] = toCelcius(line)
		}
	}
	thermal.Data = thermalMap
	return thermal
}
