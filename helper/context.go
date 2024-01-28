package helper

import (
	"errors"
	adb "github.com/abccyz/goadb"
	"github.com/basiooo/andromodem/middleware"
	"net/http"
)

func GetADBClient(r *http.Request) (*adb.Adb, error) {
	ctx := r.Context()
	adbClient, ok := ctx.Value(middleware.AdbClientKey).(*adb.Adb)
	if !ok {
		return nil, errors.New("unable to retrieve adbClient")
	}
	return adbClient, nil
}
