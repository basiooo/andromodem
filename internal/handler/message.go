package handler

import (
	"errors"
	"net/http"

	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/service"
	"github.com/basiooo/andromodem/internal/util"
	"github.com/go-chi/chi/v5"
)

type MessageHandler interface {
	GetSmsInbox(http.ResponseWriter, *http.Request)
}

type MessageHandlerImpl struct {
	MessageService service.MessageService
}

func NewMessageHander(deviceService service.MessageService) MessageHandler {
	return &MessageHandlerImpl{
		MessageService: deviceService,
	}
}

func (d *MessageHandlerImpl) GetSmsInbox(writter http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	smsInbox, err := d.MessageService.GetInbox(serial)
	if err != nil && errors.Is(err, util.ErrDeviceNotFound) {
		response := model.ErrorResponse{
			Error: "Device Not Found",
		}
		util.WriteToResponseBody(writter, response, http.StatusNotFound)
		return
	}
	util.WriteToResponseBody(writter, smsInbox, http.StatusOK)
}
