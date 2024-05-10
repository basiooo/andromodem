package service

import (
	"sync"

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

func NewMessageService(adb *adb.Adb, adbCommand adbcommand.AdbCommand) MessageService {
	return &MessageServiceImpl{
		Adb:        adb,
		AdbCommand: adbCommand,
	}
}

func (d *MessageServiceImpl) GetSmsInbox(device goadb.Device) (*[]parser.SMSInbox, string, error) {
	rawSmsInbox, _ := d.AdbCommand.GetSmsInbox(device)
	smsInboxs, err := parser.NewSMSInbox(rawSmsInbox)
	method := "without root"
	if err != nil {
		rawSmsInbox, _ := d.AdbCommand.GetSmsInboxRoot(device)
		smsInboxs, err = parser.NewSMSInbox(rawSmsInbox)
		method = "root"
	}
	return smsInboxs, method, err
}

func (d *MessageServiceImpl) GetInbox(serial string) (*model.MessageSMSInbox, error) {
	var err error
	device, err := d.Adb.GetDeviceBySerial(serial)
	if err != nil {
		return nil, util.ErrDeviceNotFound
	}

	smsInbox := &model.MessageSMSInbox{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		inboxs, method, e := d.GetSmsInbox(*device)
		err = e
		if err == nil {
			smsInbox.Inboxs = inboxs
			smsInbox.Method = &method
		}
	}()
	wg.Wait()
	if err != nil {
		return nil, err
	}
	return smsInbox, nil
}
