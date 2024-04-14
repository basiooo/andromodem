package parser

import (
	"regexp"
)

type IpAddress struct {
	Ip string
}

func NewIpAddress(rawIpAddress string) *IpAddress {
	return parseIpAddress(rawIpAddress)
}

func parseIpAddress(rawIpAddress string) *IpAddress {
	ipAddress := &IpAddress{}
	re := regexp.MustCompile(`inet addr:(\d+\.\d+\.\d+\.\d+)\s+Mask`)
	ip := re.FindStringSubmatch(rawIpAddress)
	if len(ip) > 1 {
		ipAddress.Ip = ip[1]
	}
	return ipAddress
}
