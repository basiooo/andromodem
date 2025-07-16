package parser

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

var ipRouteRegex = regexp.MustCompile(`^(\d{1,3}(?:\.\d{1,3}){3}/\d+)\s+dev\s+(\S+).*?\bsrc\s+(\d{1,3}(?:\.\d{1,3}){3})`)

type NetworkIp struct {
	Interface string `json:"interface"`
	IP        net.IP `json:"ip"`
}

type IPRoute struct {
	NetworkIPs []NetworkIp `json:"ip_routes"`
}

func NewIPRoute() IParser {
	return &IPRoute{}
}

func (i *IPRoute) Parse(rawData string) error {
	lines := strings.Split(strings.TrimSpace(rawData), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		networkIp, err := i.process(line)
		if err != nil {
			continue
		}
		if networkIp != nil {
			i.NetworkIPs = append(i.NetworkIPs, *networkIp)
		}
	}
	return nil
}

func (i *IPRoute) process(data string) (*NetworkIp, error) {
	match := ipRouteRegex.FindStringSubmatch(data)
	if len(match) != 4 {
		return nil, fmt.Errorf("no match found in input: %q", data)
	}

	parsedIP := net.ParseIP(match[3])
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid source IP: %s", match[3])
	}

	return &NetworkIp{
		Interface: match[2],
		IP:        parsedIP,
	}, nil
}
