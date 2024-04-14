package util

import (
	"encoding/json"
	"net/http"

	"github.com/basiooo/andromodem/internal/model"
)

func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	if err != nil {
		panic(err)
	}
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}, statusCode int) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	if err != nil {
		panic(err)
	}
}

func MakeDeviceNotFoundResponse(writer http.ResponseWriter) {
	response := model.BaseResponse{
		Status:  "Failed",
		Message: "Device not found",
	}
	WriteToResponseBody(writer, response, http.StatusNotFound)
}
