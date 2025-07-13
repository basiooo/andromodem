package executor

import "github.com/basiooo/andromodem/pkg/adb_processor/command"

type IExecutor interface {
	Run(command.AdbCommand) (string, error)
	EnableRoot()
	DisableRoot()
	Root() bool
}
