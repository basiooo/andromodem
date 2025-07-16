package parser_test

import (
	"net"
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseIpRoutes(t *testing.T) {
	t.Parallel()
	data := `10.43.30.144/31 dev rmnet_data2 proto kernel scope link src 10.43.30.144`
	expected := &parser.IPRoute{
		NetworkIPs: []parser.NetworkIp{
			{
				Interface: "rmnet_data2",
				IP:        net.ParseIP("10.43.30.144"),
			},
		},
	}
	ipRoute := parser.NewIPRoute()
	err := ipRoute.Parse(data)
	assert.NoError(t, err)
	ip := ipRoute.(*parser.IPRoute)
	assert.Equal(t, expected, ip)
}

func TestParseIpRoutesOlderAndroid(t *testing.T) {
	t.Parallel()
	data := `default via 10.0.2.2 dev eth0
10.0.2.0/24 dev eth0  proto kernel  scope link  src 10.0.2.15`
	expected := &parser.IPRoute{
		NetworkIPs: []parser.NetworkIp{
			{
				Interface: "eth0",
				IP:        net.ParseIP("10.0.2.15"),
			},
		},
	}
	ipRoute := parser.NewIPRoute()
	err := ipRoute.Parse(data)
	assert.NoError(t, err)
	ip := ipRoute.(*parser.IPRoute)
	assert.Equal(t, expected, ip)
}
func TestParseIpRoutesOlderAndroidVersion(t *testing.T) {
	t.Parallel()
	data := `192.168.200.0/24 dev radio0  proto kernel  scope link  src 192.168.200.2`
	expected := &parser.IPRoute{
		NetworkIPs: []parser.NetworkIp{
			{
				Interface: "radio0",
				IP:        net.ParseIP("192.168.200.2"),
			},
		},
	}
	ipRoute := parser.NewIPRoute()
	err := ipRoute.Parse(data)
	assert.NoError(t, err)
	ip := ipRoute.(*parser.IPRoute)
	assert.Equal(t, expected, ip)
}

func TestParseIpRoutesMultiple(t *testing.T) {
	t.Parallel()
	data := `10.165.252.1/30 dev rmnet_data1 proto kernel scope link src 10.165.252.241
	192.168.0.0/24 dev wlan0 proto kernel scope link src 192.168.0.157 `
	expected := &parser.IPRoute{
		NetworkIPs: []parser.NetworkIp{
			{
				Interface: "rmnet_data1",
				IP:        net.ParseIP("10.165.252.241"),
			},
			{
				Interface: "wlan0",
				IP:        net.ParseIP("192.168.0.157"),
			},
		},
	}
	ipRoute := parser.NewIPRoute()
	err := ipRoute.Parse(data)
	assert.NoError(t, err)
	ip := ipRoute.(*parser.IPRoute)
	assert.Equal(t, expected, ip)
}

func TestParseIpRoutesInvalid(t *testing.T) {
	t.Parallel()
	data := `invalid123`
	expected := &parser.IPRoute{}
	ipRoute := parser.NewIPRoute()
	err := ipRoute.Parse(data)
	assert.NoError(t, err)
	ip := ipRoute.(*parser.IPRoute)
	assert.Equal(t, expected, ip)
}

func BenchmarkParseIpRoutes(b *testing.B) {
	data := `10.43.30.144/31 dev rmnet_data2 proto kernel scope link src 10.43.30.144`
	for i := 0; i < b.N; i++ {
		ipRoute := parser.NewIPRoute()
		_ = ipRoute.Parse(data)
	}
}

func BenchmarkParseIpRoutesMultiple(b *testing.B) {
	data := `10.165.252.1/30 dev rmnet_data1 proto kernel scope link src 10.165.252.241
	192.168.0.0/24 dev wlan0 proto kernel scope link src 192.168.0.157 `
	for i := 0; i < b.N; i++ {
		ipRoute := parser.NewIPRoute()
		_ = ipRoute.Parse(data)
	}
}

func BenchmarkParseIpRoutesInvalid(b *testing.B) {
	data := `invalid123`
	for i := 0; i < b.N; i++ {
		ipRoute := parser.NewIPRoute()
		_ = ipRoute.Parse(data)
	}
}
