package server

import (
	"net/http"

	"github.com/basiooo/andromodem/internal/adb"
)

func StartServer() error {
	adb := adb.New()
	adb.Start()
	router := NewRouter(adb)
	err := http.ListenAndServe(":3001", router.Setup())
	return err
}
