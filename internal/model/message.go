package model

import "github.com/basiooo/andromodem/pkg/adb/parser"

type MessageSMSInbox struct {
	Inboxs []parser.SMSInbox `json:"inboxs"`
}
