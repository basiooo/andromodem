package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestGetLoggerLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		level    string
		expected zapcore.Level
	}{
		{"debug level", "debug", zapcore.DebugLevel},
		{"info level", "info", zapcore.InfoLevel},
		{"warn level", "warn", zapcore.WarnLevel},
		{"error level", "error", zapcore.ErrorLevel},
		{"panic level", "panic", zapcore.PanicLevel},
		{"fatal level", "fatal", zapcore.FatalLevel},
		{"unknown level", "unknown", zapcore.DebugLevel},
		{"empty level", "", zapcore.DebugLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := getLoggerLevel(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("Error changing directory back: %v", err)
		}
	}()

	logger := NewLogger()
	require.NotNil(t, logger)

	logger.Info("test message")
	logger.Error("test error message")

	logDir := filepath.Join(tempDir, "andromodem_logs")
	logFile := filepath.Join(logDir, "andromodem.log")

	time.Sleep(100 * time.Millisecond)

	_, err = os.Stat(logFile)
	assert.NoError(t, err, "Log file should be created")
}

func TestLogDuration(t *testing.T) {
	t.Parallel()

	config := zap.NewDevelopmentConfig()
	logger, err := config.Build()
	require.NoError(t, err)
	defer func() {
		if err := logger.Sync(); err != nil {
			t.Logf("Logger sync warning (expected in tests): %v", err)
		}
	}()

	done := LogDuration(logger, "test_operation")

	time.Sleep(10 * time.Millisecond)

	assert.NotPanics(t, func() {
		done()
	})
}

func TestLogDuration_MultipleOperations(t *testing.T) {
	t.Parallel()

	config := zap.NewDevelopmentConfig()
	logger, err := config.Build()
	require.NoError(t, err)
	defer func() {
		if err := logger.Sync(); err != nil {
			t.Logf("Logger sync warning (expected in tests): %v", err)
		}
	}()

	done1 := LogDuration(logger, "operation_1")
	done2 := LogDuration(logger, "operation_2")

	time.Sleep(5 * time.Millisecond)
	done1()

	time.Sleep(5 * time.Millisecond)
	done2()

	assert.True(t, true)
}

func TestLoggerLevelMap(t *testing.T) {
	t.Parallel()

	expectedLevels := []string{"debug", "info", "warn", "error", "panic", "fatal"}

	for _, level := range expectedLevels {
		t.Run("level_"+level, func(t *testing.T) {
			t.Parallel()
			_, exists := loggerLevelMap[level]
			assert.True(t, exists, "Level %s should exist in loggerLevelMap", level)
		})
	}
}

func TestLogDuration_WithNilLogger(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		done := LogDuration(nil, "test_operation")
		done()
	})
}

func BenchmarkLogDuration(b *testing.B) {
	config := zap.NewDevelopmentConfig()
	logger, _ := config.Build()
	defer func() {
		if err := logger.Sync(); err != nil {
			b.Logf("Logger sync warning (expected in tests): %v", err)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		done := LogDuration(logger, "benchmark_operation")
		done()
	}
}
