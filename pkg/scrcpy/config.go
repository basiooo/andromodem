package scrcpy

import (
	_ "embed"
)

//go:embed server/scrcpy-server-v3.3.1
var scrcpyServerBytes []byte

type Config struct {
	LocalServerBytes []byte
	RemoteServerPath string
	TCPPort          uint16
	ServerVersion    string

	Options Options
}

type Options struct {
	MaxSize uint16
	MaxFps  uint8
	Bitrate uint32
}

func NewDefaultConfigWithOptions(options *Options) *Config {
	return &Config{
		LocalServerBytes: scrcpyServerBytes,
		RemoteServerPath: SCRCPY_REMOTE_SERVER_PATH,
		TCPPort:          SCRCPY_SERVER_TCP_PORT,
		ServerVersion:    SCRCPY_SERVER_VERSION,
		Options:          *options,
	}
}
