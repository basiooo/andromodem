package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/basiooo/andromodem/internal/common"
	"github.com/basiooo/andromodem/internal/model"
	"go.uber.org/zap"
)

func Recoverer(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(writer http.ResponseWriter, request *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					func(writer http.ResponseWriter, request *http.Request, err any) {
						logger.Error(fmt.Sprintf("Panic : %v", err), zap.String("traceback", string(debug.Stack())))
						webResponse := model.BaseResponse{
							Success: false,
							Message: "Internal Server Error",
						}
						common.WriteToResponseBody(writer, webResponse, http.StatusInternalServerError)
					}(writer, request, rvr)
				}
			}()
			next.ServeHTTP(writer, request)
		}
		return http.HandlerFunc(fn)
	}

}
