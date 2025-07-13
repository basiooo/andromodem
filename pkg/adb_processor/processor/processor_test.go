package processor_test

import (
	"testing"

	"github.com/basiooo/andromodem/pkg/adb_processor/command"
	adbError "github.com/basiooo/andromodem/pkg/adb_processor/errors"
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	adbproccesor "github.com/basiooo/andromodem/pkg/adb_processor/processor"
	adb "github.com/basiooo/goadb"
	"github.com/basiooo/goadb/wire"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestProcessorWithoutDevice(t *testing.T) {
	logger := zaptest.NewLogger(t)
	processor := adbproccesor.NewProcessor(logger)
	_, err := processor.Run(nil, command.GetDeviceUptimeCommand, false)
	assert.Error(t, err)
}

func TestProcessor(t *testing.T) {
	s := &adb.MockServer{
		Status:   wire.StatusSuccess,
		Messages: []string{"10"},
	}
	device := (&adb.Adb{Server: s}).Device(adb.AnyDevice())
	logger := zaptest.NewLogger(t)
	processor := adbproccesor.NewProcessor(logger)
	result, err := processor.Run(device, command.GetAndroidVersionCommand, false)
	assert.NoError(t, err)
	androidVersion := result.(*parser.RawParser)
	assert.Equal(t, "10", androidVersion.Result)
}

func TestProcessorUnknown(t *testing.T) {
	s := &adb.MockServer{
		Status:   wire.StatusSuccess,
		Messages: []string{"boom"},
	}
	device := (&adb.Adb{Server: s}).Device(adb.AnyDevice())
	logger := zaptest.NewLogger(t)
	processor := adbproccesor.NewProcessor(logger)
	_, err := processor.Run(device, "boooom", false)
	assert.Error(t, err)
}

func TestProcessorNeedSuperUserPermission(t *testing.T) {
	s := &adb.MockServer{
		Status:   wire.StatusSuccess,
		Messages: []string{`permission denied`},
	}
	device := (&adb.Adb{Server: s}).Device(adb.AnyDevice())
	logger := zaptest.NewLogger(t)
	processor := adbproccesor.NewProcessor(logger)
	_, err := processor.Run(device, command.GetApnCommand, false)
	assert.ErrorIs(t, err, adbError.ErrorNeedShellSuperUserPermission)
}

func TestProcessorNeedRoot(t *testing.T) {
	s := &adb.MockServer{
		Status: wire.StatusSuccess,
		Messages: []string{`Error while accessing provider:telephony
	java.lang.SecurityException
		at android.os.Parcel.createException(Parcel.java:2071)
		at android.os.Parcel.readException(Parcel.java:2039)
		at android.database.DatabaseUtils.readExceptionFromParcel(DatabaseUtils.java:188)
		at android.database.DatabaseUtils.readExceptionFromParcel(DatabaseUtils.java:140)
		at android.content.ContentProviderProxy.query(ContentProviderNative.java:423)
		at com.android.commands.content.Content$QueryCommand.onExecute(Content.java:619)
		at com.android.commands.content.Content$Command.execute(Content.java:469)
		at com.android.commands.content.Content.main(Content.java:690)
		at com.android.internal.os.RuntimeInit.nativeFinishInit(Native Method)
		at com.android.internal.os.RuntimeInit.main(RuntimeInit.java:342)`},
	}
	device := (&adb.Adb{Server: s}).Device(adb.AnyDevice())
	logger := zaptest.NewLogger(t)
	processor := adbproccesor.NewProcessor(logger)
	_, err := processor.Run(device, command.GetApnCommand, false)
	assert.ErrorIs(t, err, adbError.ErrorNeedRoot)
}
