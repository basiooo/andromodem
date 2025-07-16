package messages_service

import (
	"context"

	andromodemError "github.com/basiooo/andromodem/internal/errors"
	"github.com/basiooo/andromodem/internal/service/common_service"
	"github.com/basiooo/andromodem/pkg/adb_processor/command"
	adbErrors "github.com/basiooo/andromodem/pkg/adb_processor/errors"
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
	"github.com/basiooo/andromodem/pkg/adb_processor/processor"
	"github.com/basiooo/andromodem/pkg/logger"

	adb "github.com/basiooo/goadb"
	"go.uber.org/zap"
)

type MessagesService struct {
	Adb          *adb.Adb
	AdbProcessor processor.IProcessor
	Logger       *zap.Logger
	Ctx          context.Context
}

func NewMessagesService(adb *adb.Adb, adbProcessor processor.IProcessor, logger *zap.Logger, ctx context.Context) IMessagesService {
	return &MessagesService{
		Adb:          adb,
		AdbProcessor: adbProcessor,
		Logger:       logger,
		Ctx:          ctx,
	}
}

func (m *MessagesService) GetMessages(serial string) (*parser.Inbox, error) {
	defer logger.LogDuration(m.Logger, "GetMessages")()

	device, err := m.Adb.GetDeviceBySerial(serial)
	if err != nil || device == nil {
		m.Logger.Error("error getting device by serial",
			zap.String("serial", serial),
			zap.Error(err),
		)
		return nil, andromodemError.ErrorDeviceNotFound
	}

	needRoot := false
	if androidVersion, err := common_service.GetAndroidVersion(device, m.AdbProcessor, true); err == nil {
		needRoot = androidVersion < command.MinimumAndroidShowMessages
	}

	var inbox parser.IParser

	if needRoot {
		rootInfo, err := common_service.GetDeviceRootAndAccessInfo(device, m.AdbProcessor, false)
		if err != nil {
			m.Logger.Error("inbox need root failed get root info",
				zap.String("serial", serial),
				zap.Error(err),
			)
			return nil, err
		}
		if rootInfo == nil || !rootInfo.Rooted || !rootInfo.ShellAccess {
			m.Logger.Error("inbox need root but root info is invalid or insufficient access",
				zap.String("serial", serial),
			)
			return nil, adbErrors.ErrorNeedRoot
		}

		inbox, _ = m.AdbProcessor.RunWithRoot(device, command.GetInboxCommand)
	} else {
		inbox, err = m.AdbProcessor.Run(device, command.GetInboxCommand, true)
		if err != nil {
			m.Logger.Error("error parsing inbox",
				zap.String("serial", serial),
				zap.Error(err),
			)
			return nil, err
		}
	}

	return inbox.(*parser.Inbox), nil
}
