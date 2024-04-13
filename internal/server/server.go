package server

import (
	"net/http"

	"github.com/basiooo/andromodem/internal/adb"
)

func StartServer() error {
	adbClient := adb.NewAdbClient()
	err := adbClient.Start()
	if err != nil {
		return err
	}
	router := NewRouter(adbClient)
	err = http.ListenAndServe(":3000", router.Setup())
	return err
}
