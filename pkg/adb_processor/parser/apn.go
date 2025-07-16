package parser

import (
	"encoding/json"
	"errors"
	"strings"
)

type Apn struct {
	Name     string `json:"name"`
	ApnName  string `json:"apn"`
	Proxy    string `json:"proxy"`
	Port     string `json:"port"`
	MMsProxy string `json:"mms_proxy"`
	MMsPort  string `json:"mms_port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
	MCC      string `json:"mcc"`
	MNC      string `json:"mnc"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
}

func NewApn() IParser {
	return &Apn{}
}

func (a *Apn) Parse(rawData string) error {
	parsed, err := parseApn(rawData)
	if err != nil {
		return err
	}

	if name, ok := parsed["name"]; ok {
		a.Name = name
	}
	if apn, ok := parsed["apn"]; ok {
		a.ApnName = apn
	}
	if proxy, ok := parsed["proxy"]; ok {
		a.Proxy = proxy
	}
	if port, ok := parsed["port"]; ok {
		a.Port = port
	}
	if mmsProxy, ok := parsed["mmsproxy"]; ok {
		a.MMsProxy = mmsProxy
	}
	if mmsPort, ok := parsed["mmsport"]; ok {
		a.MMsPort = mmsPort
	}
	if username, ok := parsed["user"]; ok {
		a.Username = username
	}
	if password, ok := parsed["password"]; ok {
		a.Password = password
	}
	if server, ok := parsed["server"]; ok {
		a.Server = server
	}
	if mcc, ok := parsed["mcc"]; ok {
		a.MCC = mcc
	}
	if mnc, ok := parsed["mnc"]; ok {
		a.MNC = mnc
	}
	if typ, ok := parsed["type"]; ok {
		a.Type = typ
	}
	if protocol, ok := parsed["protocol"]; ok {
		a.Protocol = protocol
	}

	return nil
}
func parseApn(rawData string) (map[string]string, error) {
	if rawData == "" {
		return nil, errors.New("empty input")
	}

	parts := strings.Split(rawData, ", ")
	result := make(map[string]string, len(parts))

	for _, part := range parts {
		if eq := strings.IndexByte(part, '='); eq != -1 {
			key := strings.TrimSpace(part[:eq])
			value := strings.TrimSpace(part[eq+1:])
			if key != "" {
				result[key] = value
			}
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no valid APN data found")
	}
	return result, nil
}
func (a Apn) MarshalJSON() ([]byte, error) {
	if a.Name == "" &&
		a.ApnName == "" &&
		a.Proxy == "" &&
		a.Port == "" &&
		a.MMsProxy == "" &&
		a.MMsPort == "" &&
		a.Username == "" &&
		a.Password == "" &&
		a.Server == "" &&
		a.MCC == "" &&
		a.MNC == "" &&
		a.Type == "" &&
		a.Protocol == "" {
		return []byte("null"), nil
	}

	type Alias Apn
	return json.Marshal(Alias(a))
}
