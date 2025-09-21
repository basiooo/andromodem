package scrcpy

import (
	"encoding/binary"
	"fmt"
	"io"

	"go.uber.org/zap"
)

type VideoStreamHandler struct {
	client  *Client
	isReady bool
}

func NewVideoStreamHandler(client *Client) *VideoStreamHandler {
	return &VideoStreamHandler{
		client: client,
	}
}

func (v *VideoStreamHandler) Reset() {
	v.isReady = false
}
func (v *VideoStreamHandler) ReadDeviceInfoHeader() error {

	if v.client.VideoConn == nil {
		return fmt.Errorf("video connection not established")
	}
	deviceNameLen := make([]byte, DEVICE_NAME_LENGTH)
	testLen := make([]byte, DUMMY_LENGTH)
	videoHeaderLen := make([]byte, VIDEO_HEADER_LENGTH)

	if _, err := io.ReadFull(v.client.VideoConn, deviceNameLen); err != nil {
		return fmt.Errorf("read device name length: %w", err)
	}

	if _, err := io.ReadFull(v.client.VideoConn, testLen); err != nil {
		return fmt.Errorf("read video test length: %w", err)
	}

	if _, err := io.ReadFull(v.client.VideoConn, videoHeaderLen); err != nil {
		return fmt.Errorf("read video header length: %w", err)
	}

	v.client.Width = binary.BigEndian.Uint32(videoHeaderLen[4:8])
	v.client.Height = binary.BigEndian.Uint32(videoHeaderLen[8:12])
	v.client.Logger.Info("Device video stream dimensions received",
		zap.Uint32("width", v.client.Width),
		zap.Uint32("height", v.client.Height))

	v.isReady = true
	return nil
}

func (v *VideoStreamHandler) CaptureVideoStream(handleChunk func([]byte)) error {
	if v.client.VideoConn == nil {
		return fmt.Errorf("video connection not established")
	}

	buf := make([]byte, 1024*1024*5) // 5mb buffer

	for {
		select {
		case <-v.client.Ctx.Done():
			if v.client.VideoConn != nil {
				if err := v.client.VideoConn.Close(); err != nil {
					v.client.Logger.Error("Failed to close video connection", zap.Error(err))
				}
			}
			return nil
		default:
			if !v.isReady {
				continue
			}
			bytesLen, err := v.client.VideoConn.Read(buf)
			if err != nil {
				return fmt.Errorf("read video stream: %w", err)
			}
			handleChunk(buf[:bytesLen])
		}
	}
}
