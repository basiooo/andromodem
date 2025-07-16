package middleware

import (
	"net/http"
	"os/exec"

	"github.com/basiooo/andromodem/internal/common"
	"github.com/basiooo/andromodem/internal/model"
	adb "github.com/basiooo/goadb"
	"go.uber.org/zap"
)

func AdbChecker(adb *adb.Adb, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(writer http.ResponseWriter, request *http.Request) {
			if adb == nil {
				logger.Error("ADB not installed in your system")
				errMessage := "ADB is not installed. Please ensure that ADB is installed on your system before proceeding."
				cmd := exec.Command("which", "adb")
				err := cmd.Run()
				if err == nil {
					errMessage = "The ADB server is currently not running on your computer. Please restart AndroModem to activate the ADB server and enable the use of this program."
				}
				webResponse := model.BaseResponse{
					Success: false,
					Message: errMessage,
				}
				common.WriteToResponseBody(writer, webResponse, http.StatusServiceUnavailable)
				return
			}
			next.ServeHTTP(writer, request)
		}
		return http.HandlerFunc(fn)
	}
}
