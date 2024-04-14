package parser_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseIpAddress(t *testing.T) {
	data := `rmnet_data1 Link encap:UNSPEC  
	inet addr:10.107.13.111  Mask:255.255.255.224 
	UP RUNNING  MTU:1500  Metric:1
	RX packets:1170 errors:0 dropped:0 overruns:0 frame:0 
	TX packets:2805 errors:0 dropped:0 overruns:0 carrier:0 
	collisions:0 txqueuelen:1000 
	RX bytes:323911 TX bytes:453957`
	expected := parser.IpAddress{
		Ip: "10.107.13.111",
	}
	ipAddress := parser.NewIpAddress(data)
	actual := *ipAddress
	assert.Equal(t, expected, actual)
}

func TestParseIpAddressNotHasIpV4(t *testing.T) {
	data := `rmnet_data1 Link encap:UNSPEC  
	inet6 addr: fe80::928:d6c0:2fbf:feea/64 Scope: Link
	UP RUNNING  MTU:1500  Metric:1
	RX packets:1170 errors:0 dropped:0 overruns:0 frame:0 
	TX packets:2805 errors:0 dropped:0 overruns:0 carrier:0 
	collisions:0 txqueuelen:1000 
	RX bytes:323911 TX bytes:453957`
	expected := parser.IpAddress{
		Ip: "",
	}
	ipAddress := parser.NewIpAddress(data)
	actual := *ipAddress
	assert.Equal(t, expected, actual)
}
