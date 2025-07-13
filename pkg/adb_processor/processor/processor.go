// Package processor provides the core logic for processing ADB commands,
// including command parsing, execution, and error handling.
package processor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/basiooo/andromodem/pkg/adb_processor/command"
	adbErrors "github.com/basiooo/andromodem/pkg/adb_processor/errors"
	"github.com/basiooo/andromodem/pkg/adb_processor/executor"
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	adb "github.com/basiooo/goadb"
	"go.uber.org/zap"
)

type Processor struct {
	Logger *zap.Logger
}

func NewProcessor(logger *zap.Logger) IProcessor {
	return &Processor{
		Logger: logger,
	}
}

// GetParser returns the appropriate parser for the given ADB command.
// Logs parser function name for debug.
func (p *Processor) GetParser(adbCommand command.AdbCommand) (parser.IParser, error) {
	fAny, ok := commandParser[adbCommand]
	if !ok {
		err := fmt.Errorf("no parser found for command: '%s'", adbCommand)
		p.Logger.Error(err.Error())
		return nil, err
	}

	f, ok := fAny.(func() parser.IParser)
	if !ok {
		err := fmt.Errorf("parser for command '%s' is not a function", adbCommand)
		p.Logger.Error(err.Error())
		return nil, err
	}

	parserInstance := f()
	p.Logger.Debug("parsing command using parser", zap.String("parser", fmt.Sprintf("%T", parserInstance)), zap.String("command", string(adbCommand)))
	return parserInstance, nil
}

// isAvailable inspects the command output for permission-related errors.
func (p *Processor) isAvailable(result string) error {
	lower := strings.ToLower(result)
	switch {
	case strings.Contains(lower, "permission denial:"),
		strings.Contains(lower, "permission denied"):
		return adbErrors.ErrorNeedShellSuperUserPermission
	case strings.Contains(lower, "error while accessing provider"),
		strings.Contains(lower, "su: invalid uid/gid"):
		return adbErrors.ErrorNeedRoot
	default:
		return nil
	}
}

// runWithPermissionCheck runs a command and checks for permission errors in output.
func (p *Processor) runWithPermissionCheck(exec executor.IExecutor, adbCommand command.AdbCommand) (string, error) {
	result, err := exec.Run(adbCommand)
	if err != nil {
		return "", err
	}

	if permErr := p.isAvailable(result); permErr != nil {
		return "", permErr
	}

	return result, nil
}

// runCommand is a helper method to execute ADB commands and process the results.
// This method extracts common logic from Run and RunWithRoot.
func (p *Processor) runCommand(device *adb.Device, adbCommand command.AdbCommand, useRoot bool) (parser.IParser, error) {
	if device == nil {
		return nil, adbErrors.ErrorDeviceIsNil
	}

	exec := executor.NewExecutor(device)

	parserInstance, err := p.GetParser(adbCommand)
	if err != nil {
		return nil, err
	}

	if useRoot {
		exec.EnableRoot()
	}

	result, err := p.runWithPermissionCheck(exec, adbCommand)
	if err != nil {
		p.Logger.Error("failed to execute command", zap.Error(err), zap.Bool("with_root", useRoot))
		return nil, err
	}

	if err := parserInstance.Parse(result); err != nil {
		p.Logger.Error("failed to parse result", zap.Error(err))
		return nil, err
	}

	return parserInstance, nil
}

// Run executes an ADB command on the given device, optionally using root if permission denied.
// Returns the parsed result or an error.
func (p *Processor) Run(device *adb.Device, adbCommand command.AdbCommand, userRootIfDenied bool) (parser.IParser, error) {
	result, err := p.runCommand(device, adbCommand, false)
	if err != nil && userRootIfDenied && (errors.Is(err, adbErrors.ErrorNeedShellSuperUserPermission) || errors.Is(err, adbErrors.ErrorNeedRoot)) {
		// Try again with root if permission denied and userRootIfDenied is true
		return p.runCommand(device, adbCommand, true)
	}
	return result, err
}

// RunWithRoot executes an ADB command on the given device with root privileges enabled.
// Returns the parsed result or an error.
func (p *Processor) RunWithRoot(device *adb.Device, adbCommand command.AdbCommand) (parser.IParser, error) {
	return p.runCommand(device, adbCommand, true)
}
