package middleware

import (
	"net/http"

	"github.com/basiooo/andromodem/internal/adb"
	"github.com/basiooo/andromodem/internal/model"
	"github.com/basiooo/andromodem/internal/util"
)

func AdbChecker(adb *adb.Adb) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(writer http.ResponseWriter, request *http.Request) {
			if adb.Client == nil {
				errMessage := "ADB is not installed. Please ensure that ADB is installed on your system before proceeding."
				if adb.AdbIsInstalled() {
					errMessage = "The ADB server is currently not running on your computer. Please restart Andromodem to activate the ADB server and enable the use of this program."
				}
				webResponse := model.BaseResponse{
					Status:  "Failed",
					Message: errMessage,
				}
				util.WriteToResponseBody(writer, webResponse, http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(writer, request)
		}
		return http.HandlerFunc(fn)
	}
}
