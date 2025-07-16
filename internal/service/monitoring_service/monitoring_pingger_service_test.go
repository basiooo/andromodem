package monitoring_service

import (
	"testing"

	adb "github.com/basiooo/goadb"
	"go.uber.org/zap"
)

func TestMonitoringPinggerService_parsePingResult(t *testing.T) {
	t.Parallel()
	adb := &adb.Adb{}
	logger := zap.NewNop()
	service := NewMonitoringPinggerService(adb, logger).(*MonitoringPinggerService)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Whitespace only",
			input:    "   ",
			expected: false,
		},
		{
			name:     "Standard Linux ping success",
			input:    "PING google.com (142.250.191.14) 56(84) bytes of data.\n64 bytes from lga25s62-in-f14.1e100.net (142.250.191.14): icmp_seq=1 ttl=118 time=12.3 ms\n\n--- google.com ping statistics ---\n1 packets transmitted, 1 received, 0% packet loss, time 0ms\nrtt min/avg/max/mdev = 12.345/12.345/12.345/0.000 ms",
			expected: true,
		},
		{
			name:     "MacOS ping success",
			input:    "PING google.com (142.250.191.14): 56 data bytes\n64 bytes from 142.250.191.14: icmp_seq=0 ttl=118 time=12.345 ms\n\n--- google.com ping statistics ---\n1 packets transmitted, 1 received, 0% packet loss\nround-trip min/avg/max/stddev = 12.345/12.345/12.345/0.000 ms",
			expected: true,
		},
		{
			name:     "Windows ping success",
			input:    "Pinging google.com [142.250.191.14] with 32 bytes of data:\nReply from 142.250.191.14: bytes=32 time=12ms TTL=118\n\nPing statistics for 142.250.191.14:\n    Packets: Sent = 1, Received = 1, Lost = 0 (0% loss)",
			expected: true,
		},
		{
			name:     "Case insensitive - uppercase",
			input:    "1 PACKETS TRANSMITTED, 1 RECEIVED, 0% PACKET LOSS",
			expected: true,
		},
		{
			name:     "Case insensitive - mixed case",
			input:    "1 packets transmitted, 1 Received, 0% packet loss",
			expected: true,
		},
		{
			name:     "Multiple packets with 1 received",
			input:    "5 packets transmitted, 1 received, 80% packet loss",
			expected: true,
		},
		{
			name:     "Zero percent packet loss",
			input:    "5 packets transmitted, 5 received, 0% packet loss",
			expected: true,
		},
		{
			name:     "Zero percent with decimal",
			input:    "10 packets transmitted, 10 received, 0.0% packet loss",
			expected: true,
		},
		{
			name:     "Zero percent uppercase",
			input:    "1 PACKETS TRANSMITTED, 1 RECEIVED, 0% PACKET LOSS",
			expected: true,
		},
		{
			name:     "No packets received",
			input:    "1 packets transmitted, 0 received, 100% packet loss",
			expected: false,
		},
		{
			name:     "Partial packet loss",
			input:    "5 packets transmitted, 3 received, 40% packet loss",
			expected: false,
		},
		{
			name:     "High packet loss",
			input:    "10 packets transmitted, 2 received, 80% packet loss",
			expected: false,
		},
		{
			name:     "Network unreachable",
			input:    "ping: connect: Network is unreachable",
			expected: false,
		},
		{
			name:     "Host unreachable",
			input:    "From 192.168.1.1 icmp_seq=1 Destination Host Unreachable",
			expected: false,
		},
		{
			name:     "Timeout",
			input:    "Request timeout for icmp_seq 0",
			expected: false,
		},
		{
			name:     "Multiple lines with success",
			input:    "Some error occurred\nBut later: 1 packets transmitted, 1 received, 0% packet loss\nEnd of output",
			expected: true,
		},
		{
			name:     "Contains 1 received but not in context",
			input:    "Error: 1 received error message, connection failed",
			expected: true,
		},
		{
			name:     "Contains 0% packet loss but not in ping context",
			input:    "System report: 0% packet loss in network interface statistics",
			expected: true,
		},
		{
			name:     "Numbers without context",
			input:    "Random text with numbers 1 2 3 received some data",
			expected: false,
		},
		{
			name:     "Packet loss but not 0%",
			input:    "1 packets transmitted, 0 received, 100% packet loss",
			expected: false,
		},
		{
			name:     "Android ping success",
			input:    "PING 8.8.8.8 (8.8.8.8) 56(84) bytes of data.\n64 bytes from 8.8.8.8: icmp_seq=1 ttl=118 time=25.4 ms\n\n--- 8.8.8.8 ping statistics ---\n1 packets transmitted, 1 received, 0% packet loss, time 0ms",
			expected: true,
		},
		{
			name:     "Android ping failure",
			input:    "PING 192.168.999.999 (192.168.999.999) 56(84) bytes of data.\n\n--- 192.168.999.999 ping statistics ---\n1 packets transmitted, 0 received, 100% packet loss, time 0ms",
			expected: false,
		},
		{
			name:     "With special characters",
			input:    "\t\n1 packets transmitted, 1 received, 0% packet loss\r\n",
			expected: true,
		},
		{
			name:     "Very long output with success",
			input:    "Very long ping output with lots of details and information about network configuration and routing tables and finally 1 packets transmitted, 1 received, 0% packet loss at the end",
			expected: true,
		},
		{
			name:     "Malformed output",
			input:    "ping: invalid option -- 'z'\nTry 'ping --help' for more information.",
			expected: false,
		},
		{
			name:     "Permission denied",
			input:    "ping: socket: Operation not permitted",
			expected: false,
		},
		{
			name:     "DNS resolution failure",
			input:    "ping: cannot resolve invalid-domain-name.com: Unknown host",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.isIMCPSuccess(tt.input)
			if result != tt.expected {
				t.Errorf("isIMCPSuccess() = %v, expected %v\nInput: %q", result, tt.expected, tt.input)
			}
		})
	}
}
