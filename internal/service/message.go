package service

import (
	"fmt"

	goadb "github.com/abccyz/goadb"
	"github.com/basiooo/andromodem/internal/adb"
	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/util"
	adbcommand "github.com/basiooo/andromodem/pkg/adb/adb_command"
	"github.com/basiooo/andromodem/pkg/adb/parser"
)

type MessageService interface {
	GetInbox(string) (*model.MessageSMSInbox, error)
}

type MessageServiceImpl struct {
	*adb.Adb
	AdbCommand adbcommand.AdbCommand
}

func NewMessageService(adbClient *adb.Adb, adbCommand adbcommand.AdbCommand) MessageService {
	return &MessageServiceImpl{
		Adb:        adbClient,
		AdbCommand: adbCommand,
	}
}
func (d *MessageServiceImpl) GetSmsInbox(device goadb.Device) []parser.SMSInbox {
	rawSmsInbox, _ := d.AdbCommand.GetSmsInbox(device)
	fmt.Println(rawSmsInbox)
	smsInboxs := parser.NewSMSInbox(rawSmsInbox)
	return smsInboxs
}

func (d *MessageServiceImpl) GetInbox(serial string) (*model.MessageSMSInbox, error) {
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil {
		return nil, util.ErrDeviceNotFound
	}
	smsInbox := &model.MessageSMSInbox{}
	smsInboxChan := make(chan []parser.SMSInbox)
	go func() {
		defer close(smsInboxChan)
		smsInboxChan <- d.GetSmsInbox(*device)
	}()
	smsInbox.Inboxs = <-smsInboxChan

	return smsInbox, nil
}
