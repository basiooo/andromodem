package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/basiooo/andromodem/internal/router"
	"github.com/basiooo/andromodem/internal/server"
	"github.com/basiooo/andromodem/internal/utils"
	"github.com/basiooo/andromodem/pkg/cache"
	"github.com/basiooo/andromodem/pkg/logger"
	adb "github.com/basiooo/goadb"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var Version = "dev"

func main() {
	versionFlag := flag.Bool("version", false, "Show version information")
	vFlag := flag.Bool("v", false, "Show version information")
	flag.Parse()

	if *versionFlag || *vFlag {
		fmt.Printf("AndroModem version %s\n", Version)
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Create log directories
	logDirs := []string{
		"andromodem_logs",
		"andromodem_logs/monitoring",
	}
	for _, dir := range logDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Failed to create log directory %s: %v\n", dir, err)
		}
	}

	appLogger := logger.NewLogger()
	appLogger.Info("Application starting", zap.String("version", Version))

	_ = cache.NewCache(5*time.Minute, 10*time.Minute)

	adbClient, err := adb.New()
	if err != nil {
		appLogger.Error("Failed create ADB client", zap.String("error", err.Error()))
	}
	validator := validator.New()
	validator.RegisterTagNameFunc(utils.GetJSONFieldName)
	router := router.NewRouter(adbClient, appLogger, ctx, validator)
	server := server.NewServer(router.GetRouters(), appLogger, ctx, cancel)
	if err := server.Start(); err != nil {
		fmt.Println(err.Error())
	}
}
