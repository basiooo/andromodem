package scrcpy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"time"

	adb "github.com/basiooo/goadb"

	"go.uber.org/zap"
)

type Client struct {
	Device      *adb.Device
	VideoConn   net.Conn
	ControlConn net.Conn
	Config      *Config
	Width       uint32
	Height      uint32
	Logger      *zap.Logger
	Ctx         context.Context

	IsRunning bool

	VideoHandler   *VideoStreamHandler
	ControlHandler *ControlHandler
}

func NewClient(ctx context.Context, device *adb.Device, config *Config, logger *zap.Logger) *Client {
	client := &Client{
		Device: device,
		Ctx:    ctx,
		Config: config,
		Logger: logger,
	}
	client.VideoHandler = NewVideoStreamHandler(client)
	client.ControlHandler = NewControlHandler(client)
	return client
}

func (c *Client) SetupPortForward() error {
	if err := c.Device.ForwardAbstract(c.Config.TCPPort, "scrcpy"); err != nil {
		c.Logger.Error("Failed to set up port forwarding", zap.Error(err))
		return fmt.Errorf("fail setup port forwarding: %w", err)
	}
	c.Logger.Debug("Port forwarding set up", zap.Uint16("port", c.Config.TCPPort))
	return nil
}

func (c *Client) PushServer() error {
	c.Logger.Debug("Pushing scrcpy server to device")
	localFile := bytes.NewReader(c.Config.LocalServerBytes)
	remoteWriter, err := c.Device.OpenWrite(c.Config.RemoteServerPath, 0644, time.Now())
	if err != nil {
		c.Logger.Error("Failed to open remote file", zap.Error(err))
		return fmt.Errorf("fail open remote file: %w", err)
	}
	defer func() {
		if err := remoteWriter.Close(); err != nil {
			c.Logger.Error("Failed to close remote file", zap.Error(err))
		}
	}()

	if _, err := io.Copy(remoteWriter, localFile); err != nil {
		c.Logger.Error("Failed to push scrcpy server", zap.Error(err))
		return fmt.Errorf("fail push server: %w", err)
	}
	c.Logger.Info("Scrcpy server pushed to device")
	return nil
}

func (c *Client) StartServer() error {
	c.Logger.Info("Starting scrcpy server")
	if err := c.PushServer(); err != nil {
		return err
	}

	cmd := fmt.Sprintf(
		"CLASSPATH=%s app_process / com.genymobile.scrcpy.Server %s tunnel_forward=true audio=false control=true send_frame_meta=false send_device_meta=true max_size=%d max_fps=%d video_bit_rate=%d power_on=true",
		c.Config.RemoteServerPath,
		c.Config.ServerVersion,
		c.Config.Options.MaxSize,
		c.Config.Options.MaxFps,
		c.Config.Options.Bitrate,
	)
	c.Logger.Debug("scrcpy server command", zap.String("cmd", cmd))

	go func() {
		err := c.Device.RunShellLoop(c.Ctx, cmd)
		if err != nil {
			if adb.HasErrCode(err, adb.CtxCanceled) {
				c.Logger.Info("Scrcpy server loop canceled")
			} else {
				c.Logger.Error("Scrcpy server exited with error", zap.Error(err))
			}
		} else {
			c.Logger.Info("Scrcpy server loop exited cleanly")
		}
	}()

	// Wait for server startup
	// TODO: Replace hardcoded sleep with proper server readiness check
	// Current implementation uses a fixed 3-second delay which is unreliable.
	time.Sleep(3 * time.Second)
	c.Logger.Debug("Assumed scrcpy server started (3s delay)")
	return nil
}

func (c *Client) ResetState() {
	c.Logger.Info("Resetting scrcpy client state")
	c.Width = 0
	c.Height = 0
	c.IsRunning = false
	c.VideoHandler.Reset()
}

func (c *Client) Cleanup() error {
	c.Logger.Info("Cleaning up scrcpy client resources")
	c.ResetState()

	if err := c.Device.ForwardRemove(fmt.Sprintf("tcp:%d", c.Config.TCPPort)); err != nil {
		c.Logger.Warn("Failed to remove port forwarding", zap.Error(err))
		return fmt.Errorf("remove port forwarding: %w", err)
	}

	c.Logger.Info("Scrcpy client cleanup complete")
	return nil
}

func (c *Client) ConnectToStream() error {
	videoConn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", c.Config.TCPPort))
	if err != nil {
		c.Logger.Error("Failed to connect to video stream", zap.Error(err))
		return fmt.Errorf("fail connect to video stream: %w", err)
	}
	c.VideoConn = videoConn

	controlConn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", c.Config.TCPPort))
	if err != nil {
		c.Logger.Error("Failed to connect to control stream", zap.Error(err))
		return fmt.Errorf("fail connect to control stream: %w", err)
	}
	c.ControlConn = controlConn

	if tcpConn, ok := c.VideoConn.(*net.TCPConn); ok {
		if err := tcpConn.SetNoDelay(true); err != nil {
			c.Logger.Error("Failed to set TCP no delay", zap.Error(err))
			return fmt.Errorf("set TCP no delay: %w", err)
		}
		if err := tcpConn.SetKeepAlive(true); err != nil {
			c.Logger.Error("Failed to set TCP keep alive", zap.Error(err))
			return fmt.Errorf("set TCP keep alive: %w", err)
		}
		if err := tcpConn.SetKeepAlivePeriod(30 * time.Second); err != nil {
			c.Logger.Error("Failed to set TCP keep alive period", zap.Error(err))
			return fmt.Errorf("set TCP keep alive period: %w", err)
		}
	}

	if tcpConn, ok := c.ControlConn.(*net.TCPConn); ok {
		if err := tcpConn.SetNoDelay(true); err != nil {
			c.Logger.Error("Failed to set TCP no delay", zap.Error(err))
			return fmt.Errorf("set TCP no delay: %w", err)
		}
		if err := tcpConn.SetKeepAlive(true); err != nil {
			c.Logger.Error("Failed to set TCP keep alive", zap.Error(err))
			return fmt.Errorf("set TCP keep alive: %w", err)
		}
		if err := tcpConn.SetKeepAlivePeriod(30 * time.Second); err != nil {
			c.Logger.Error("Failed to set TCP keep alive period", zap.Error(err))
			return fmt.Errorf("set TCP keep alive period: %w", err)
		}
	}

	c.Logger.Info("Successfully connected to video and control streams")
	return nil
}

func (c *Client) Start() error {
	c.Logger.Info("Initializing scrcpy client session")
	if err := c.SetupPortForward(); err != nil {
		return err
	}
	if err := c.StartServer(); err != nil {
		return err
	}
	if err := c.ConnectToStream(); err != nil {
		return err
	}
	if err := c.VideoHandler.ReadDeviceInfoHeader(); err != nil {
		return fmt.Errorf("read device info header: %w", err)
	}

	go func() {
		<-c.Ctx.Done()
		c.Logger.Warn("Context canceled, cleaning up scrcpy client")
		if err := c.Cleanup(); err != nil {
			c.Logger.Error("Cleanup failed", zap.Error(err))
		}
	}()

	return nil
}
