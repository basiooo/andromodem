package model

import "github.com/basiooo/andromodem/pkg/adb/parser"

type MessageSMSInbox struct {
	Method *string            `json:"method"`
	Inboxs *[]parser.SMSInbox `json:"inboxs"`
}
