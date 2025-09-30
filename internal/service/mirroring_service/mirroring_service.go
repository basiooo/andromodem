package mirroring_service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	andromodemError "github.com/basiooo/andromodem/internal/errors"
	"github.com/basiooo/andromodem/pkg/scrcpy"
	adb "github.com/basiooo/goadb"
	"go.uber.org/zap"
)

type MirroringService struct {
	Adb     *adb.Adb
	Logger  *zap.Logger
	Ctx     context.Context
	clients map[string]*scrcpy.Client

	mutex sync.RWMutex
}

func NewMirroringService(adb *adb.Adb, logger *zap.Logger, ctx context.Context) IMirroringService {
	return &MirroringService{
		Adb:     adb,
		Logger:  logger,
		Ctx:     ctx,
		clients: make(map[string]*scrcpy.Client),
	}
}

func (m *MirroringService) StartMirroring(ctx context.Context, serial string) (*scrcpy.Client, error) {
	device, err := m.Adb.GetDeviceBySerial(serial)
	if err != nil || device == nil {
		m.Logger.Error("error getting device by serial",
			zap.String("serial", serial),
			zap.Error(err),
		)
		return nil, andromodemError.ErrorDeviceNotFound
	}
	m.mutex.RLock()
	if client, exists := m.clients[serial]; exists {
		m.mutex.RUnlock()
		m.Logger.Info("mirroring already started for device", zap.String("serial", serial))
		return client, nil
	}
	m.mutex.RUnlock()

	config := scrcpy.NewDefaultConfigWithOptions(&scrcpy.Options{
		MaxSize: 1080,
		MaxFps:  60,
	})

	client := scrcpy.NewClient(ctx, device, config, m.Logger)

	if err := client.Start(); err != nil {
		m.Logger.Error("failed to start scrcpy client",
			zap.String("serial", serial),
			zap.Error(err),
		)
		return nil, fmt.Errorf("start scrcpy client: %w", err)
	}

	m.mutex.Lock()
	m.clients[serial] = client
	m.mutex.Unlock()

	m.Logger.Info("scrcpy mirroring started successfully",
		zap.String("serial", serial),
	)
	go func() {
		<-ctx.Done()
		m.mutex.Lock()
		defer m.mutex.Unlock()
		if _, exists := m.clients[serial]; exists {
			delete(m.clients, serial)
			m.Logger.Info("Stopping scrcpy mirroring",
				zap.String("serial", serial),
			)
		}
	}()
	return client, nil
}

func (m *MirroringService) CaptureVideoStream(serial string, handleVideoChunk func([]byte)) error {
	m.mutex.RLock()
	client, exists := m.clients[serial]
	m.mutex.RUnlock()
	if !exists {
		return fmt.Errorf("mirroring not started for device %s", serial)
	}
	m.Logger.Info("starting video stream capture", zap.String("serial", serial))
	if err := client.VideoHandler.CaptureVideoStream(handleVideoChunk); err != nil {
		m.Logger.Error("error capturing video stream",
			zap.String("serial", serial),
			zap.Error(err),
		)
		return fmt.Errorf("capture video stream: %w", err)
	}

	return nil
}

func (m *MirroringService) IsRunning(serial string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	_, exists := m.clients[serial]
	return exists
}

func (m *MirroringService) SendTouchEvent(serial string, event *scrcpy.TouchEvent) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	client, exists := m.clients[serial]

	if !exists {
		return fmt.Errorf("mirroring not started for device %s", serial)
	}

	if err := client.ControlHandler.SendTouchEvent(event); err != nil {
		m.Logger.Error("failed to send touch event",
			zap.String("serial", serial),
			zap.Error(err),
		)
		return fmt.Errorf("send touch event: %w", err)
	}
	return nil
}

func (m *MirroringService) SendKeyEvent(serial string, event *scrcpy.KeyEvent) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	client, exists := m.clients[serial]
	if !exists {
		return fmt.Errorf("mirroring not started for device %s", serial)
	}

	if err := client.ControlHandler.SendKeyEvent(event); err != nil {
		m.Logger.Error("failed to send key event",
			zap.String("serial", serial),
			zap.Error(err),
		)
		return fmt.Errorf("send key event: %w", err)
	}
	return nil
}

func (m *MirroringService) SendKeyPress(serial string, keyCode scrcpy.AndroidKeyCode) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	client, exists := m.clients[serial]
	if !exists {
		return fmt.Errorf("mirroring not started for device %s", serial)
	}

	if err := client.ControlHandler.SendKeyPress(keyCode); err != nil {
		m.Logger.Error("failed to send key press",
			zap.String("serial", serial),
			zap.Error(err),
		)
		return fmt.Errorf("send key press: %w", err)
	}
	return nil
}

func (m *MirroringService) HandleControlMessage(serial string, messageBytes []byte) {
	if len(messageBytes) == 0 {
		m.Logger.Warn("received empty control message", zap.String("serial", serial))
		return
	}

	if len(messageBytes) > 512 {
		m.Logger.Warn("control message too large, ignoring",
			zap.String("serial", serial),
			zap.Int("size", len(messageBytes)),
		)
		return
	}

	var message map[string]interface{}
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		m.Logger.Error("failed to parse control message",
			zap.String("serial", serial),
			zap.Error(err),
			zap.String("raw_message", string(messageBytes)),
		)
		return
	}

	msgType, ok := message["type"].(string)
	if !ok {
		m.Logger.Warn("control message missing type field",
			zap.String("serial", serial),
		)
		return
	}

	allowedTypes := map[string]bool{
		"touch":  true,
		"key":    true,
		"scroll": true,
	}
	if msgType == "ping" {
		m.Logger.Debug("received ping message", zap.String("serial", serial))
		return
	}

	if !allowedTypes[msgType] {
		m.Logger.Warn("unknown or disallowed control message type",
			zap.String("serial", serial),
			zap.String("type", msgType),
		)
		return
	}

	switch msgType {
	case "touch":
		m.HandleTouchMessage(serial, message)
	case "key":
		m.HandleKeyMessage(serial, message)
	default:
		m.Logger.Warn("unhandled control message type",
			zap.String("serial", serial),
			zap.String("type", msgType),
		)
	}
}

func (m *MirroringService) HandleKeyMessage(serial string, message map[string]any) {
	key, ok := message["key"].(string)
	if !ok {
		m.Logger.Warn("key message missing key field", zap.String("serial", serial))
		return
	}

	keyCommand := scrcpy.KeyCommand(key)
	keyCode, exists := scrcpy.KeyCodeMap[keyCommand]
	if !exists {
		m.Logger.Warn("unsupported key",
			zap.String("serial", serial),
			zap.String("key", string(keyCommand)),
		)
		return
	}

	m.Logger.Info("handling key press",
		zap.String("serial", serial),
		zap.String("key", string(keyCommand)),
		zap.Int32("keycode", int32(keyCode)),
	)

	if err := m.SendKeyPress(serial, keyCode); err != nil {
		m.Logger.Error("failed to send key press",
			zap.String("serial", serial),
			zap.String("key", string(keyCommand)),
			zap.Error(err),
		)
	}
}

func (m *MirroringService) HandleTouchMessage(serial string, message map[string]any) {
	touchMsg := &scrcpy.TouchMessage{}

	touchMsg.Type = scrcpy.MessageTypeTouch

	action, ok := message["action"].(string)
	if !ok {
		m.Logger.Warn("touch message missing action", zap.String("serial", serial))
		return
	}

	touchAction := scrcpy.TouchMessageAction(action)

	allowedActions := map[scrcpy.TouchMessageAction]bool{
		scrcpy.TouchMessageActionDown:   true,
		scrcpy.TouchMessageActionUp:     true,
		scrcpy.TouchMessageActionMove:   true,
		scrcpy.TouchMessageActionCancel: true,
	}

	if !allowedActions[touchAction] {
		m.Logger.Warn("invalid touch action",
			zap.String("serial", serial),
			zap.String("action", string(touchAction)),
		)
		return
	}
	touchMsg.Action = touchAction

	x, ok := message["x"].(float64)
	if !ok {
		m.Logger.Warn("touch message missing x coordinate", zap.String("serial", serial))
		return
	}
	if x < 0.0 || x > 1.0 {
		m.Logger.Warn("invalid x coordinate",
			zap.String("serial", serial),
			zap.Float64("x", x),
		)
		return
	}
	touchMsg.X = float32(x)

	y, ok := message["y"].(float64)
	if !ok {
		m.Logger.Warn("touch message missing y coordinate", zap.String("serial", serial))
		return
	}
	if y < 0.0 || y > 1.0 {
		m.Logger.Warn("invalid y coordinate",
			zap.String("serial", serial),
			zap.Float64("y", y),
		)
		return
	}
	touchMsg.Y = float32(y)

	if pointerId, ok := message["pointerId"].(float64); ok {
		if pointerId < 0 || pointerId > 10 {
			m.Logger.Warn("invalid pointer ID",
				zap.String("serial", serial),
				zap.Float64("pointerId", pointerId),
			)
			return
		}
		touchMsg.PointerId = uint64(pointerId)
	} else {
		touchMsg.PointerId = 0
	}

	if pressure, ok := message["pressure"].(float64); ok {
		if pressure < 0.0 || pressure > 1.0 {
			m.Logger.Warn("invalid pressure value",
				zap.String("serial", serial),
				zap.Float64("pressure", pressure),
			)
			return
		}
		touchMsg.Pressure = float32(pressure)
	} else {
		touchMsg.Pressure = 1.0
	}

	m.Logger.Debug("processing touch message",
		zap.String("serial", serial),
		zap.String("action", string(touchMsg.Action)),
		zap.Float32("x", touchMsg.X),
		zap.Float32("y", touchMsg.Y),
		zap.Uint64("pointerId", touchMsg.PointerId),
		zap.Float32("pressure", touchMsg.Pressure),
	)

	touchEvent := scrcpy.ConvertTouchMessage(touchMsg)

	if err := m.SendTouchEvent(serial, touchEvent); err != nil {
		m.Logger.Error("failed to send touch event",
			zap.String("serial", serial),
			zap.Error(err),
			zap.String("action", string(touchMsg.Action)),
			zap.Float32("x", touchMsg.X),
			zap.Float32("y", touchMsg.Y),
		)
		return
	}

	m.Logger.Debug("touch event processed successfully",
		zap.String("serial", serial),
		zap.String("action", string(touchMsg.Action)),
		zap.Float32("x", touchMsg.X),
		zap.Float32("y", touchMsg.Y),
	)
}
func (m *MirroringService) GetClient(serial string) *scrcpy.Client {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.clients[serial]
}
