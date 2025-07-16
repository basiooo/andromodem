package executor

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/command"
	adb "github.com/basiooo/goadb"
	"github.com/basiooo/goadb/wire"
	"github.com/stretchr/testify/assert"
)

func mockDevice() (*adb.MockServer, *adb.Device) {
	s := &adb.MockServer{
		Status:   wire.StatusSuccess,
		Messages: []string{"output"},
	}
	return s, (&adb.Adb{Server: s}).Device(adb.AnyDevice())
}

func TestCommandEnableDisable(t *testing.T) {
	t.Parallel()
	_, device := mockDevice()
	executor := NewExecutor(device)
	executor.EnableRoot()
	assert.True(t, executor.Root(), "Root should be true after calling EnableRoot")

	executor.DisableRoot()
	assert.False(t, executor.Root(), "Root should be false after calling DisableRoot")
}

func TestCommandWithRoot(t *testing.T) {
	t.Parallel()
	s, device := mockDevice()
	executor := NewExecutor(device)
	executor.EnableRoot()
	output1, err1 := executor.Run(command.AdbCommand("test"))
	assert.NoError(t, err1)
	assert.Equal(t, "shell:su -c 'test'", s.Requests[1])
	assert.Equal(t, "output", output1)
}

func TestCommandWithoutRoot(t *testing.T) {
	t.Parallel()
	s, device := mockDevice()
	executor := NewExecutor(device)
	executor.DisableRoot()
	output, err := executor.Run(command.AdbCommand("test"))
	assert.NoError(t, err)
	assert.Equal(t, "shell:test", s.Requests[1])
	assert.Equal(t, "output", output)
}
