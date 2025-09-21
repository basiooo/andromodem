package mirroring_service

import (
	"context"

	"github.com/basiooo/andromodem/pkg/scrcpy"
)

type IMirroringService interface {
	StartMirroring(context.Context, string) (*scrcpy.Client, error)
	CaptureVideoStream(string, func([]byte)) error
	IsRunning(string) bool
	SendTouchEvent(string, *scrcpy.TouchEvent) error
	SendKeyEvent(string, *scrcpy.KeyEvent) error
	SendKeyPress(string, scrcpy.AndroidKeyCode) error
	HandleControlMessage(string, []byte)
	GetClient(string) *scrcpy.Client
}
