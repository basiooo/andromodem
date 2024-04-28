package parser

import (
	"net"
	"strings"
)

type IpAddress struct {
	Ip string
}

func NewIpAddress(rawIpAddress string) *IpAddress {
	return parseIpAddress(rawIpAddress)
}

func parseIpAddress(rawIpAddress string) *IpAddress {
	ipAddress := &IpAddress{}
	if net.IP(rawIpAddress) != nil {
		ipAddress.Ip = strings.TrimSpace(rawIpAddress)
	}
	return ipAddress
}
