package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

type SignalStrength struct {
	*CellSignalStrengthCdma    `json:"Cdma,omitempty"`
	*CellSignalStrengthGsm     `json:"Gsm,omitempty"`
	*CellSignalStrengthWcdma   `json:"Wcdma,omitempty"`
	*CellSignalStrengthTdscdma `json:"Tdscdma,omitempty"`
	*CellSignalStrengthLte     `json:"Lte,omitempty"`
	*CellSignalStrengthNr      `json:"Nr,omitempty"`
}

func (ss *SignalStrength) IsEmpty() bool {
	return ss.CellSignalStrengthCdma == nil && ss.CellSignalStrengthGsm == nil && ss.CellSignalStrengthWcdma == nil && ss.CellSignalStrengthTdscdma == nil && ss.CellSignalStrengthLte == nil && ss.CellSignalStrengthNr == nil
}

type CellSignalStrengthCdma struct {
	Level    int    `json:"level,omitempty"`
	CdmaDbm  string `json:"cdmaDbm,omitempty"`
	CdmaEcio string `json:"cdmaEcio,omitempty"`
	EvdoDbm  string `json:"evdoDbm,omitempty"`
	EvdoEcio string `json:"evdoEcio,omitempty"`
	EvdoSnr  string `json:"evdoSnr,omitempty"`
}

type CellSignalStrengthGsm struct {
	Level string `json:"mLevel,omitempty"`
	Rssi  string `json:"rssi,omitempty"`
	Ber   string `json:"ber,omitempty"`
	Ta    string `json:"mTa,omitempty"`
}

type CellSignalStrengthWcdma struct {
	Level string `json:"level,omitempty"`
	Ss    string `json:"ss,omitempty"`
	Ber   string `json:"ber,omitempty"`
	Rscp  string `json:"rscp,omitempty"`
	Ecno  string `json:"ecno,omitempty"`
}

type CellSignalStrengthTdscdma struct {
	Level string `json:"level,omitempty"`
	Rssi  string `json:"rssi,omitempty"`
	Ber   string `json:"ber,omitempty"`
	Rscp  string `json:"rscp,omitempty"`
}

type CellSignalStrengthLte struct {
	Level                 string `json:"level,omitempty"`
	Rssi                  string `json:"rssi,omitempty"`
	Rsrp                  string `json:"rsrp,omitempty"`
	Rsrq                  string `json:"rsrq,omitempty"`
	Rssnr                 string `json:"rssnr,omitempty"`
	CqiTableIndex         string `json:"cqiTableIndex,omitempty"`
	Cqi                   string `json:"cqi,omitempty"`
	Ta                    string `json:"ta,omitempty"`
	ParametersUseForLevel string `json:"parametersUseForLevel,omitempty"`
}

type CellSignalStrengthNr struct {
	Level                 string   `json:"level,omitempty"`
	CsiRsrp               string   `json:"csiRsrp,omitempty"`
	CsiRsrq               string   `json:"csiRsrq,omitempty"`
	CsiCqiTableIndex      string   `json:"csiCqiTableIndex,omitempty"`
	CsiCqiReport          []string `json:"csiCqiReport,omitempty"`
	SsRsrp                string   `json:"ssRsrp,omitempty"`
	SsRsrq                string   `json:"ssRsrq,omitempty"`
	SsSinr                string   `json:"ssSinr,omitempty"`
	ParametersUseForLevel string   `json:"parametersUseForLevel,omitempty"`
	TimingAdvance         string   `json:"timingAdvance,omitempty"`
}

func NewSignalStrength(rawsignalStrength string) *SignalStrength {
	signalStrength := parseSignalStrength(rawsignalStrength)
	if signalStrength.IsEmpty() {
		signalStrength = parseSignalStrengthOld(rawsignalStrength)
	}
	return signalStrength
}

func parseSignalStrengthOld(rawsignalStrength string) *SignalStrength {
	var networkType string
	regexNetworkType := regexp.MustCompile(`\s\w+\|\w+`)
	signalStrength := &SignalStrength{}
	match := regexNetworkType.FindString(rawsignalStrength)
	networkTypes := strings.Split(match, "|")
	if len(networkTypes) == 1 {
		networkType = networkTypes[0]
	} else if len(networkTypes) == 2 {
		networkType = networkTypes[1]
	} else {
		return signalStrength
	}
	networkType = strings.TrimSpace(networkType)
	networkType = strings.ToLower(networkType)
	switch networkType {
	case "cdma":
		signalStrength.CellSignalStrengthCdma = &CellSignalStrengthCdma{}
	case "gsm":
		signalStrength.CellSignalStrengthGsm = &CellSignalStrengthGsm{}
	case "wcdma":
		signalStrength.CellSignalStrengthWcdma = &CellSignalStrengthWcdma{}
	case "tdscdma":
		signalStrength.CellSignalStrengthTdscdma = &CellSignalStrengthTdscdma{}
	case "lte":
		signalStrength.CellSignalStrengthLte = &CellSignalStrengthLte{}
	case "nr":
		signalStrength.CellSignalStrengthNr = &CellSignalStrengthNr{}
	default:

	}
	return signalStrength
}

func parseSignalStrength(rawsignalStrength string) *SignalStrength {
	regexNetworkType := regexp.MustCompile(`primary=([^,]+)(?:,|})`)
	match := regexNetworkType.FindStringSubmatch(rawsignalStrength)
	if len(match) < 2 {
		return &SignalStrength{}
	}
	networkType := strings.TrimSpace(match[1])
	signalStrengthJson := extractSignalStrength(networkType, rawsignalStrength)
	signalStrength := createCellSignalStrength(networkType, signalStrengthJson)
	return signalStrength
}

func extractSignalStrength(networkType string, rawSignalStrength string) string {
	regexSignalStrengthValue := regexp.MustCompile(fmt.Sprintf(`%s: ([^,]+),`, networkType))
	value := regexSignalStrengthValue.FindStringSubmatch(rawSignalStrength)
	if len(value) < 2 {
		return ""
	}
	return formatRawSignalStrength(value[1])

}

func createCellSignalStrength(networkType, signalStrengthJson string) *SignalStrength {
	cellSignalStrength := &SignalStrength{}

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
		logrus.WithField("location", "NewCarriersInfo.createCellSignalStrength").Error(err)
		return cellSignalStrength
	}

	return cellSignalStrength

}

func formatRawSignalStrength(data string) string {
	pairs := strings.Split(data, " ")
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
		logrus.WithField("location", "NewCarriersInfo.formatRawSignalStrength").Error(err)
		return ""
	}
	return string(jsonData)
}
