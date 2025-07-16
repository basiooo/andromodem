package rest

import (
	"errors"
	"net/http"

	"github.com/basiooo/andromodem/internal/common"
	andromodemError "github.com/basiooo/andromodem/internal/errors"
	"github.com/basiooo/andromodem/internal/service/messages_service"
	adbError "github.com/basiooo/andromodem/pkg/adb_processor/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type MessagesHandler struct {
	MessageService messages_service.IMessagesService
	Logger         *zap.Logger
	Validator      *validator.Validate
}

func NewMessagesHandler(messageService messages_service.IMessagesService, logger *zap.Logger, validator *validator.Validate) IMessagesHandler {
	return &MessagesHandler{
		MessageService: messageService,
		Logger:         logger,
		Validator:      validator,
	}
}

func (m MessagesHandler) GetMessages(writer http.ResponseWriter, request *http.Request) {
	serial := chi.URLParam(request, "serial")
	messages, err := m.MessageService.GetMessages(serial)
	if err != nil {
		if errors.Is(err, andromodemError.ErrorDeviceNotFound) {
			common.DeviceNotFoundResponse(writer)
			return
		}
		if errors.Is(err, adbError.ErrorMinimumAndroidVersionNotSupport) || errors.Is(err, adbError.ErrorNeedShellSuperUserPermission) || errors.Is(err, adbError.ErrorNeedRoot) {
			common.ErrorResponse(writer, err.Error(), http.StatusServiceUnavailable)
			return
		}
		common.ErrorResponse(writer, "Error getting messages", http.StatusInternalServerError)
		return
	}
	common.SuccessResponse(writer, "messages retrieved successfully", messages, http.StatusOK)
}
