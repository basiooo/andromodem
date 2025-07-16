package processor

import (
	"github.com/basiooo/andromodem/pkg/adb_processor/command"
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	adb "github.com/basiooo/goadb"
)

type IProcessor interface {
	GetParser(command.AdbCommand) (parser.IParser, error)
	Run(*adb.Device, command.AdbCommand, bool) (parser.IParser, error)
	RunWithRoot(*adb.Device, command.AdbCommand) (parser.IParser, error)
}
