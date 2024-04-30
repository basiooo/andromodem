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
	response := model.BaseResponse{
		Status:  "Success",
		Message: "SMS list retrieved successfully",
	}
	statusCode := http.StatusOK
	if err != nil {
		if errors.Is(err, util.ErrDeviceNotFound) {
			util.MakeDeviceNotFoundResponse(writter)
			return
		}
		response.Status = "Failed"
		response.Message = err.Error()
		statusCode = http.StatusInternalServerError
	}
	response.Data = smsInbox
	util.WriteToResponseBody(writter, response, statusCode)
}
