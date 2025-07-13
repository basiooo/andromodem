package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type SignalStrength struct {
	*CellSignalStrengthCdma    `json:"Cdma,omitempty"`
	*CellSignalStrengthGsm     `json:"Gsm,omitempty"`
	*CellSignalStrengthWcdma   `json:"Wcdma,omitempty"`
	*CellSignalStrengthTdscdma `json:"Tdscdma,omitempty"`
	*CellSignalStrengthLte     `json:"Lte,omitempty"`
	*CellSignalStrengthNr      `json:"Nr,omitempty"`
}

func (s *SignalStrength) IsEmpty() bool {
	return s.CellSignalStrengthCdma == nil &&
		s.CellSignalStrengthGsm == nil &&
		s.CellSignalStrengthWcdma == nil &&
		s.CellSignalStrengthTdscdma == nil &&
		s.CellSignalStrengthLte == nil &&
		s.CellSignalStrengthNr == nil
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

var (
	regNetworkType          = regexp.MustCompile(`primary=([^,]+)(?:,|})`)
	regexSignalStrengthBase = `(?m)%s: ([^,]+),`
)

func NewSignalStrength() IParser {
	return &SignalStrength{}
}

func (s *SignalStrength) Parse(rawData string) error {
	rawData = strings.TrimSpace(rawData)
	if rawData == "" {
		return nil
	}

	err := s.parseSignalStrength(rawData)
	if s.IsEmpty() {
		err = s.parseSignalStrengthOlderAndroidVersion(strings.ToLower(rawData))
	}

	return err
}

func (s *SignalStrength) parseSignalStrengthOlderAndroidVersion(rawData string) error {
	switch rawData {
	case "hscsd", "gsm", "gprs", "edge":
		s.CellSignalStrengthGsm = &CellSignalStrengthGsm{}
	case "umts", "hsdpa":
		s.CellSignalStrengthWcdma = &CellSignalStrengthWcdma{}
	case "lte":
		s.CellSignalStrengthLte = &CellSignalStrengthLte{}
	default:
		return errors.New("unsupported older android signal type")
	}
	return nil
}

func (s *SignalStrength) parseSignalStrength(rawData string) error {
	match := regNetworkType.FindStringSubmatch(rawData)
	if len(match) < 2 {
		return errors.New("failed to parse signal strength: network type not found")
	}
	networkType := strings.TrimSpace(match[1])

	signalStrengthJson := extractSignalStrength(networkType, rawData)
	if signalStrengthJson == "" {
		return errors.New("failed to extract signal strength data for network type " + networkType)
	}

	signalMap := map[string]any{
		"CellSignalStrengthCdma":    &s.CellSignalStrengthCdma,
		"CellSignalStrengthGsm":     &s.CellSignalStrengthGsm,
		"CellSignalStrengthWcdma":   &s.CellSignalStrengthWcdma,
		"CellSignalStrengthTdscdma": &s.CellSignalStrengthTdscdma,
		"CellSignalStrengthLte":     &s.CellSignalStrengthLte,
		"CellSignalStrengthNr":      &s.CellSignalStrengthNr,
	}

	ptr, ok := signalMap[networkType]
	if !ok {
		return fmt.Errorf("unknown network type: %s", networkType)
	}

	if err := json.Unmarshal([]byte(signalStrengthJson), ptr); err != nil {
		return fmt.Errorf("failed to unmarshal signal strength JSON: %w", err)
	}

	return nil
}

func extractSignalStrength(networkType, rawData string) string {
	regex := regexp.MustCompile(fmt.Sprintf(regexSignalStrengthBase, regexp.QuoteMeta(networkType)))
	matches := regex.FindStringSubmatch(rawData)
	if len(matches) < 2 {
		return ""
	}
	return formatRawSignalStrength(matches[1])
}

func formatRawSignalStrength(data string) string {
	if strings.TrimSpace(data) == "" {
		return "{}"
	}
	pairs := strings.Fields(data)
	result := make(map[string]string, len(pairs))

	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}
