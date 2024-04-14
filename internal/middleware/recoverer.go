package middleware

import (
	"net/http"

	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/util"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				func(writer http.ResponseWriter, request *http.Request, err interface{}) {
					errorMessage := "Internal Server Error"
					if message, ok := err.(string); ok {
						errorMessage = message
					}
					webResponse := model.BaseResponse{
						Status:  "Failed",
						Message: errorMessage,
					}
					util.WriteToResponseBody(writer, webResponse, http.StatusInternalServerError)
				}(writer, request, rvr)
			}
		}()
		next.ServeHTTP(writer, request)
	}

	return http.HandlerFunc(fn)
}
