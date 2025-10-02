package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/service/mirroring_service"
	"github.com/basiooo/andromodem/pkg/scrcpy"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type MirroringHandler struct {
	MirroringService mirroring_service.IMirroringService
	upgrader         websocket.Upgrader
	Logger           *zap.Logger
	Validator        *validator.Validate

	activeConns map[string]*websocket.Conn
	mu          sync.Mutex
}

func NewMirroringHandler(mirroringService mirroring_service.IMirroringService, logger *zap.Logger, validator *validator.Validate) *MirroringHandler {
	return &MirroringHandler{
		MirroringService: mirroringService,
		upgrader: websocket.Upgrader{
			CheckOrigin:      func(r *http.Request) bool { return true },
			HandshakeTimeout: 10 * time.Second,
		},
		Logger:      logger,
		activeConns: make(map[string]*websocket.Conn),
		Validator:   validator,
	}
}

func (h *MirroringHandler) StartMirroringStream(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")
	var client *scrcpy.Client
	var err error
	conn, err := h.upgrader.Upgrade(w, r, nil)
	h.mu.Lock()
	if oldConn, ok := h.activeConns[serial]; ok {
		h.Logger.Info("closing old websocket connection", zap.String("serial", serial))
		if err := oldConn.Close(); err != nil {
			h.Logger.Error("failed to close old websocket connection",
				zap.String("serial", serial),
				zap.Error(err),
			)
		}
		for {
			if ok := h.MirroringService.GetClient(serial); ok == nil {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
	h.activeConns[serial] = conn
	h.mu.Unlock()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	conn.SetCloseHandler(func(code int, text string) error {
		h.Logger.Info("websocket connection closed by client",
			zap.String("serial", serial),
			zap.Int("code", code),
			zap.String("text", text),
		)
		message := websocket.FormatCloseMessage(code, "")
		if err := conn.WriteControl(
			websocket.CloseMessage,
			message,
			time.Now().Add(time.Second),
		); err != nil {
			h.Logger.Warn("failed to write close message", zap.Error(err))
		}
		cancel()
		return nil
	})

	if err != nil {
		h.Logger.Error("failed to upgrade to websocket",
			zap.String("serial", serial),
			zap.Error(err),
		)
		return
	}
	h.Logger.Info("websocket connection established for mirroring",
		zap.String("serial", serial),
	)

	setupReceived := make(chan *model.MirroringSetupRequest, 1)

	go h.handleWebSocketMessages(ctx, cancel, conn, serial, setupReceived)

	var setup *model.MirroringSetupRequest
	select {
	case setup = <-setupReceived:
		if !h.MirroringService.IsRunning(serial) {
			client, err = h.MirroringService.StartMirroring(ctx, serial, setup)
			if err != nil {
				h.Logger.Error("failed to start mirroring",
					zap.String("serial", serial),
					zap.Error(err),
				)
				errorMsg := map[string]interface{}{
					"type":    "error",
					"message": "Failed to start mirroring: " + err.Error(),
				}
				if err := conn.WriteJSON(errorMsg); err != nil {
					h.Logger.Error("failed to send error message",
						zap.String("serial", serial),
						zap.Error(err),
					)
				}
				return
			}
		}
	case <-ctx.Done():
		return
	}

	successMsg := map[string]interface{}{
		"type":    "connected",
		"message": "Mirroring stream connected",
		"serial":  serial,
		"width":   client.Width,
		"height":  client.Height,
	}
	if err := conn.WriteJSON(successMsg); err != nil {
		h.Logger.Error("failed to send success message",
			zap.String("serial", serial),
			zap.Error(err),
		)
		return
	}

	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	handleVideoChunk := func(chunk []byte) {
		if len(chunk) == 0 {
			return
		}

		if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
			h.Logger.Warn("failed to set write deadline",
				zap.String("serial", serial),
				zap.Error(err),
			)
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, chunk); err != nil {
			h.Logger.Warn("failed to send video chunk",
				zap.String("serial", serial),
				zap.Int("chunk_size", len(chunk)),
				zap.Error(err),
			)
			cancel()
			return
		}
	}
	go func() {
		if err := h.MirroringService.CaptureVideoStream(serial, handleVideoChunk); err != nil {
			h.Logger.Error("video stream capture ended",
				zap.String("serial", serial),
				zap.Error(err),
			)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			h.Logger.Info("mirroring websocket session ended",
				zap.String("serial", serial),
			)
			h.mu.Lock()
			delete(h.activeConns, serial)
			h.mu.Unlock()

			return
		case <-pingTicker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				h.Logger.Error("failed to send ping",
					zap.String("serial", serial),
					zap.Error(err),
				)
				return
			}
		}
	}
}

func (h *MirroringHandler) handleWebSocketMessages(ctx context.Context, cancelCtx context.CancelFunc, conn *websocket.Conn, serial string, setupChan chan<- *model.MirroringSetupRequest) {
	conn.SetReadLimit(1024)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				h.Logger.Info("client disconnected",
					zap.String("serial", serial),
					zap.Error(err),
				)
				cancelCtx()
				return
			}

			if len(messageBytes) > 1024 {
				h.Logger.Warn("received oversized message, ignoring",
					zap.String("serial", serial),
					zap.Int("size", len(messageBytes)),
				)
				continue
			}

			if len(messageBytes) > 0 {
				var message map[string]interface{}
				if err := json.Unmarshal(messageBytes, &message); err == nil {
					if msgType, ok := message["type"].(string); ok && msgType == "setup" {
						var setupRequest model.MirroringSetupRequest
						if err := json.Unmarshal(messageBytes, &setupRequest); err != nil {
							h.Logger.Error("failed to parse setup message",
								zap.String("serial", serial),
								zap.Error(err),
							)
							errorMsg := map[string]interface{}{
								"type":    "error",
								"message": "Invalid setup message format",
							}
							if err := conn.WriteJSON(errorMsg); err != nil {
								h.Logger.Error("failed to send error message",
									zap.String("serial", serial),
									zap.Error(err),
								)
							}
							continue
						}

						if err := h.Validator.Struct(&setupRequest); err != nil {
							h.Logger.Error("setup validation failed",
								zap.String("serial", serial),
								zap.Error(err),
							)
							errorMsg := map[string]interface{}{
								"type":    "error",
								"message": "Setup validation failed: " + err.Error(),
							}
							if err := conn.WriteJSON(errorMsg); err != nil {
								h.Logger.Error("failed to send validation error message",
									zap.String("serial", serial),
									zap.Error(err),
								)
							}
							continue
						}

						h.Logger.Info("received valid setup configuration",
							zap.String("serial", serial),
							zap.Uint8("fps", setupRequest.FPS),
							zap.Uint32("bitrate", setupRequest.Bitrate),
							zap.Uint16("resolution", setupRequest.Resolution),
						)

						select {
						case setupChan <- &setupRequest:
						default:
							h.Logger.Warn("setup channel full, ignoring duplicate setup message",
								zap.String("serial", serial),
							)
						}
						continue
					}
				}

				h.MirroringService.HandleControlMessage(serial, messageBytes)
			}
		}
	}
}
