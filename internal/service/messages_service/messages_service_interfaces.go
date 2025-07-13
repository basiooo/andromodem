package messages_service

import (
	"github.com/basiooo/andromodem/pkg/adb_processor/parser"
)

type IMessagesService interface {
	GetMessages(string) (*parser.Inbox, error)
}
