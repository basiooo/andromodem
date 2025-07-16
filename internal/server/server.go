package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

type Server interface {
	Start() error
}

type httpServer struct {
	server    *http.Server
	logger    *zap.Logger
	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewServer(router http.Handler, log *zap.Logger, ctx context.Context, ctxCancel context.CancelFunc) Server {
	return &httpServer{
		server: &http.Server{
			Addr:    "0.0.0.0:49153",
			Handler: router,
		},
		ctx:       ctx,
		ctxCancel: ctxCancel,
		logger:    log,
	}
}

func (s *httpServer) Start() error {
	s.logger.Info("Starting server")
	s.logger.Info(fmt.Sprintf("Server running on %s", s.server.Addr))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(sig)
	defer close(sig)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		select {
		case <-sig:
			s.ctxCancel()
			s.logger.Info("Shutdown signal received")
		case <-s.ctx.Done():
			s.logger.Error("Context cancelled before server started")
			return
		}

		if err := s.server.Close(); err != nil {
			s.logger.Error("Server forced to shutdown", zap.String("error", err.Error()))
		}
	}()

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("Server error", zap.String("error", err.Error()))
		return err
	}

	wg.Wait()
	s.logger.Info("Server stopped")
	return nil
}
