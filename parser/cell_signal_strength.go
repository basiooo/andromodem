package parser

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

type CellSignalStrength struct {
	*CellSignalStrengthCdma    `json:"Cdma,omitempty"`
	*CellSignalStrengthGsm     `json:"Gsm,omitempty"`
	*CellSignalStrengthWcdma   `json:"Wcdma,omitempty"`
	*CellSignalStrengthTdscdma `json:"Tdscdma,omitempty"`
	*CellSignalStrengthLte     `json:"Lte,omitempty"`
	*CellSignalStrengthNr      `json:"Nr,omitempty"`
}

type CellSignalStrengthCdma struct {
	CdmaDbm  string `json:"cdmaDbm,omitempty"`
	CdmaEcio string `json:"cdmaEcio,omitempty"`
	EvdoDbm  string `json:"evdoDbm,omitempty"`
	EvdoEcio string `json:"evdoEcio,omitempty"`
	EvdoSnr  string `json:"evdoSnr,omitempty"`
	Level    int    `json:"level,omitempty"`
}

type CellSignalStrengthGsm struct {
	Rssi  string `json:"rssi,omitempty"`
	Ber   string `json:"ber,omitempty"`
	Ta    string `json:"mTa,omitempty"`
	Level string `json:"mLevel,omitempty"`
}

type CellSignalStrengthWcdma struct {
	Ss    string `json:"ss,omitempty"`
	Ber   string `json:"ber,omitempty"`
	Rscp  string `json:"rscp,omitempty"`
	Ecno  string `json:"ecno,omitempty"`
	Level string `json:"level,omitempty"`
}

type CellSignalStrengthTdscdma struct {
	Rssi  string `json:"rssi,omitempty"`
	Ber   string `json:"ber,omitempty"`
	Rscp  string `json:"rscp,omitempty"`
	Level string `json:"level,omitempty"`
}

type CellSignalStrengthLte struct {
	Rssi                  string `json:"rssi,omitempty"`
	Rsrp                  string `json:"rsrp,omitempty"`
	Rsrq                  string `json:"rsrq,omitempty"`
	Rssnr                 string `json:"rssnr,omitempty"`
	CqiTableIndex         string `json:"cqiTableIndex,omitempty"`
	Cqi                   string `json:"cqi,omitempty"`
	Ta                    string `json:"ta,omitempty"`
	Level                 string `json:"level,omitempty"`
	ParametersUseForLevel string `json:"parametersUseForLevel,omitempty"`
}

type CellSignalStrengthNr struct {
	CsiRsrp               string   `json:"csiRsrp,omitempty"`
	CsiRsrq               string   `json:"csiRsrq,omitempty"`
	CsiCqiTableIndex      string   `json:"csiCqiTableIndex,omitempty"`
	CsiCqiReport          []string `json:"csiCqiReport,omitempty"`
	SsRsrp                string   `json:"ssRsrp,omitempty"`
	SsRsrq                string   `json:"ssRsrq,omitempty"`
	SsSinr                string   `json:"ssSinr,omitempty"`
	Level                 string   `json:"level,omitempty"`
	ParametersUseForLevel string   `json:"parametersUseForLevel,omitempty"`
	TimingAdvance         string   `json:"timingAdvance,omitempty"`
}

func formatRawSignalStrength(rawSignalStrength string) string {
	pairs := strings.Split(rawSignalStrength, " ")
	result := make(map[string]interface{})
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			result[key] = value
		}
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		logrus.WithField("function", "formatRawSignalStrength").Error(err)
		return ""
	}
	return string(jsonData)
}

func ParseSignalStrength(rawSignalStrength string) []CellSignalStrength {
	listRawSignalStrength := strings.Split(strings.TrimSpace(rawSignalStrength), "\n")
	cellSignalStrengths := make([]CellSignalStrength, 0, len(listRawSignalStrength))

	for _, signalStrengthRaw := range listRawSignalStrength {
		regexNetworkType := regexp.MustCompile(`primary=([^,]+)}`)
		match := regexNetworkType.FindStringSubmatch(signalStrengthRaw)
		if len(match) < 2 {
			continue
		}
		networkType := strings.TrimSpace(match[1])

		signalStrengthJson := extractSignalStrength(signalStrengthRaw, networkType)
		cellSignalStrength := createCellSignalStrength(networkType, signalStrengthJson)
		cellSignalStrengths = append(cellSignalStrengths, cellSignalStrength)
	}
	return cellSignalStrengths
}

func extractSignalStrength(signalStrengthRaw, networkType string) string {
	regexSignalStrengthValue := regexp.MustCompile(fmt.Sprintf(`%s: ([^,]+),`, networkType))
	value := regexSignalStrengthValue.FindStringSubmatch(signalStrengthRaw)
	if len(value) < 2 {
		return ""
	}
	return formatRawSignalStrength(value[1])
}

func createCellSignalStrength(networkType, signalStrengthJson string) CellSignalStrength {
	cellSignalStrength := CellSignalStrength{}

	signalStrengthStructs := map[string]interface{}{
		"CellSignalStrengthLte":     &cellSignalStrength.CellSignalStrengthLte,
		"CellSignalStrengthWcdma":   &cellSignalStrength.CellSignalStrengthWcdma,
		"CellSignalStrengthGsm":     &cellSignalStrength.CellSignalStrengthGsm,
		"CellSignalStrengthCdma":    &cellSignalStrength.CellSignalStrengthCdma,
		"CellSignalStrengthTdscdma": &cellSignalStrength.CellSignalStrengthTdscdma,
		"CellSignalStrengthNr":      &cellSignalStrength.CellSignalStrengthNr,
	}

	signalStrengthStruct, found := signalStrengthStructs[networkType]
	if !found {
		return cellSignalStrength
	}

	err := json.Unmarshal([]byte(signalStrengthJson), signalStrengthStruct)
	if err != nil {
		logrus.WithField("function", "createCellSignalStrength").Error(err)
		return cellSignalStrength
	}

	return cellSignalStrength
}
