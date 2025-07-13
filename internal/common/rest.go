package common

import (
	"encoding/json"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result any) error {
	decoder := json.NewDecoder(request.Body)
	return decoder.Decode(result)
}

func WriteToResponseBody(writer http.ResponseWriter, response any, statusCode int) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	if response == nil && statusCode == http.StatusNoContent {
		return
	}
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	if err != nil {
		panic(err)
	}
}
