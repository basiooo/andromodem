package utils

import (
	"os/exec"

	"go.uber.org/zap"
)

// InitializeADB initializes the ADB daemon and triggers device discovery.
// This is required on some OpenWRT devices where connected Android devices
// are not detected automatically until `adb devices` is executed.
func InitializeADB(logger *zap.Logger) error {
	logger.Info("Initializing ADB daemon and triggering device discovery...")

	cmd := exec.Command("adb", "devices")
	output, err := cmd.Output()
	if err != nil {
		logger.Error("Failed to initialize ADB daemon",
			zap.String("error", err.Error()),
			zap.String("hint", "Please install Android Debug Bridge (ADB) and ensure it's in your PATH"))
		return err
	}

	logger.Info("ADB daemon initialized successfully",
		zap.String("output", string(output)))
	return nil
}
