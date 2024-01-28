package exception

import (
	"github.com/basiooo/andromodem/helper"
	"github.com/basiooo/andromodem/model"
	"net/http"
)

func PanicHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	errorMessage := "Internal Server Error"
	if message, ok := err.(string); ok {
		errorMessage = message
	}
	webResponse := model.WebResponseError{
		Error: errorMessage,
	}

	helper.WriteToResponseBody(writer, webResponse, http.StatusInternalServerError)
}
