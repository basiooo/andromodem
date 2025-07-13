// Package executor provides functionality for executing ADB commands on a connected device.
package executor

import (
	"fmt"
	"strings"

	"github.com/basiooo/andromodem/pkg/adb_processor/command"
	adb "github.com/basiooo/goadb"
)

type Executor struct {
	Device  *adb.Device
	UseRoot bool
}

func NewExecutor(device *adb.Device) IExecutor {
	return &Executor{
		Device: device,
	}
}
func (c *Executor) Root() bool {
	return c.UseRoot
}
func (c *Executor) EnableRoot() {
	c.UseRoot = true
}
func (c *Executor) DisableRoot() {
	c.UseRoot = false
}

func (c *Executor) Run(adbCommand command.AdbCommand) (string, error) {
	cmdStr := string(adbCommand)

	if c.UseRoot {
		if !strings.HasPrefix(cmdStr, "su ") {
			cmdStr = fmt.Sprintf("su -c '%s'", cmdStr)
		}
		return c.Device.RunCommandWithTimeout(cmdStr, 10)
	}

	return c.Device.RunCommand(cmdStr)
}
