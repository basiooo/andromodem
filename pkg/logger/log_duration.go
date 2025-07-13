package logger

import (
	"time"

	"go.uber.org/zap"
)

func LogDuration(logger *zap.Logger, label string) func() {
	start := time.Now()
	return func() {
		// Check if logger is nil to prevent panic
		if logger != nil {
			logger.Debug("duration", zap.String("step", label), zap.Duration("took", time.Since(start)))
		}
		// If logger is nil, silently ignore the logging
	}
}
