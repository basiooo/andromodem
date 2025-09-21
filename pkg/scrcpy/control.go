package scrcpy

import (
	"encoding/binary"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type TouchEvent struct {
	Action AndroidTouchAction

	PointerId uint64

	X        float32
	Y        float32
	Pressure float32
	Buttons  uint32
}

type TouchMessage struct {
	Type      MessageType        `json:"type"`
	Action    TouchMessageAction `json:"action"`
	X         float32            `json:"x"`
	Y         float32            `json:"y"`
	PointerId uint64             `json:"pointerId"`
	Pressure  float32            `json:"pressure"`
}
type ControlHandler struct {
	client      *Client
	VideoBuffer []byte
}

func NewControlHandler(client *Client) *ControlHandler {
	return &ControlHandler{
		client: client,
	}
}

func (c *ControlHandler) CleanUp() {
	c.VideoBuffer = nil
}

func (c *ControlHandler) SendTouchEvent(event *TouchEvent) error {
	if c.client.ControlConn == nil {
		return fmt.Errorf("control connection not established")
	}

	if event.X < 0 || event.X > 1 || event.Y < 0 || event.Y > 1 {
		return fmt.Errorf("invalid coordinates: x=%f, y=%f (must be 0.0-1.0)", event.X, event.Y)
	}

	absX := int32(event.X * float32(c.client.Width))
	absY := int32(event.Y * float32(c.client.Height))

	buf := make([]byte, 32)
	offset := 0

	buf[offset] = byte(ControlMsgTypeInjectTouchEvent)
	offset++

	buf[offset] = byte(event.Action)
	offset++

	binary.BigEndian.PutUint64(buf[offset:], event.PointerId)
	offset += 8

	binary.BigEndian.PutUint32(buf[offset:], uint32(absX))
	offset += 4

	binary.BigEndian.PutUint32(buf[offset:], uint32(absY))
	offset += 4

	binary.BigEndian.PutUint16(buf[offset:], uint16(c.client.Width))
	offset += 2

	binary.BigEndian.PutUint16(buf[offset:], uint16(c.client.Height))
	offset += 2

	pressureInt := uint16(0xFFFF)
	if event.Pressure > 0 && event.Pressure <= 1.0 {
		pressureInt = uint16(event.Pressure * 0xFFFF)
	}
	binary.BigEndian.PutUint16(buf[offset:], pressureInt)
	offset += 2

	binary.BigEndian.PutUint32(buf[offset:], uint32(event.Buttons))
	offset += 4

	binary.BigEndian.PutUint32(buf[offset:], 0)
	offset += 4

	c.client.Logger.Debug("sending touch event",
		zap.Int32("action", int32(event.Action)),
		zap.Float32("x", event.X),
		zap.Float32("y", event.Y),
		zap.Int32("absX", absX),
		zap.Int32("absY", absY),
		zap.Uint32("width", c.client.Width),
		zap.Uint32("height", c.client.Height),
	)

	c.client.ControlConn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	n, err := c.client.ControlConn.Write(buf)
	if err != nil {
		c.client.Logger.Error("failed to send touch event",
			zap.Error(err),
			zap.Int32("action", int32(event.Action)),
			zap.Int("bytes_written", n),
		)
		return fmt.Errorf("send touch event: %w", err)
	}

	if n != 32 {
		c.client.Logger.Warn("incomplete touch event sent",
			zap.Int("expected", 32),
			zap.Int("actual", n),
		)
	}

	c.client.Logger.Debug("touch event sent successfully",
		zap.Int32("action", int32(event.Action)),
		zap.Int32("pointerId", int32(event.PointerId)),
		zap.Float32("pressure", event.Pressure),
		zap.Float32("x", event.X),
		zap.Float32("y", event.Y),
		zap.Int32("buttons", int32(event.Buttons)),
		zap.Int("bytes_written", n),
	)

	return nil
}

type KeyEvent struct {
	Action    AndroidKeyAction
	KeyCode   AndroidKeyCode
	Repeat    uint32
	MetaState uint32
}

type KeyMessage struct {
	Type MessageType `json:"type"`
	Key  KeyCommand  `json:"key"`
}

func (c *ControlHandler) SendKeyEvent(event *KeyEvent) error {
	if c.client.ControlConn == nil {
		return fmt.Errorf("control connection not established")
	}

	buf := make([]byte, 14)
	offset := 0

	buf[offset] = byte(ControlMsgTypeInjectKeyEvent)
	offset++

	buf[offset] = byte(event.Action)
	offset++

	binary.BigEndian.PutUint32(buf[offset:], uint32(event.KeyCode))
	offset += 4

	binary.BigEndian.PutUint32(buf[offset:], event.Repeat)
	offset += 4

	binary.BigEndian.PutUint32(buf[offset:], event.MetaState)
	offset += 4

	c.client.Logger.Debug("sending key event",
		zap.Int32("action", int32(event.Action)),
		zap.Int32("keycode", int32(event.KeyCode)),
		zap.Uint32("repeat", event.Repeat),
		zap.Uint32("metastate", event.MetaState),
		zap.Int("buffer_size", offset),
	)

	c.client.ControlConn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	n, err := c.client.ControlConn.Write(buf)
	if err != nil {
		c.client.Logger.Error("failed to send key event",
			zap.Error(err),
			zap.Int32("keycode", int32(event.KeyCode)),
			zap.Int("bytes_written", n),
		)
		return fmt.Errorf("send key event: %w", err)
	}

	if n != 14 {
		c.client.Logger.Warn("incomplete key event sent",
			zap.Int("expected", 14),
			zap.Int("actual", n),
		)
	}

	c.client.Logger.Debug("key event sent successfully",
		zap.Int32("action", int32(event.Action)),
		zap.Int32("keycode", int32(event.KeyCode)),
		zap.Int("bytes_written", n),
	)

	return nil
}

func (c *ControlHandler) SendKeyPress(keyCode AndroidKeyCode) error {
	if err := c.SendKeyEvent(&KeyEvent{
		Action:    AndroidKeyActionDown,
		KeyCode:   keyCode,
		Repeat:    0,
		MetaState: 0,
	}); err != nil {
		return fmt.Errorf("send key down: %w", err)
	}

	time.Sleep(50 * time.Millisecond)

	if err := c.SendKeyEvent(&KeyEvent{
		Action:    AndroidKeyActionUp,
		KeyCode:   keyCode,
		Repeat:    0,
		MetaState: 0,
	}); err != nil {
		return fmt.Errorf("send key up: %w", err)
	}

	return nil
}
