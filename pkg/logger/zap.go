package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var loggerLevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}

func getLoggerLevel(level string) zapcore.Level {
	if lvl, ok := loggerLevelMap[level]; ok {
		return lvl
	}
	return zapcore.DebugLevel
}

func NewLogger() *zap.Logger {
	logLevel := getLoggerLevel("info")
	zapConfig := zap.NewProductionEncoderConfig()
	var writeSyncer zapcore.WriteSyncer
	appDir, _ := os.Getwd()

	logOutput := filepath.Join(appDir, "andromodem_logs/andromodem.log")
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logOutput,
		MaxSize:    1,
		MaxBackups: 1,
	}

	defer func() {
		if err := lumberjackLogger.Close(); err != nil {
			fmt.Println("failed close lumberjack")
		}
	}()
	writeSyncer = zapcore.AddSync(lumberjackLogger)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConfig),
		writeSyncer,
		logLevel,
	)
	return zap.New(core)
}
