package parser

import (
	"encoding/json"
	"errors"
	"strings"
)

type Apn struct {
	Name    string `json:"name"`
	ApnName string `json:"apn"`
}

func NewApn(rawApn string) *Apn {
	apn := &Apn{}
	extrackedApn := extractApn(rawApn)
	parseApn(extrackedApn, apn)
	return apn
}

func extractApn(rawApn string) map[string]interface{} {
	result := make(map[string]interface{})
	datas := strings.Split(rawApn, ", ")
	for _, data := range datas {
		parts := strings.Split(data, "=")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result
}

func parseApn(data map[string]interface{}, apn *Apn) error {
	apnJson, err := json.Marshal(data)
	if err != nil {
		return errors.New("failed convert apn interface to json")
	}
	err = json.Unmarshal([]byte(apnJson), &apn)
	if err != nil {
		return errors.New("failed convert apn json to Battery struct")
	}
	return nil
}
