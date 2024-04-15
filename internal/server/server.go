package server

import (
	"net/http"

	"github.com/basiooo/andromodem/internal/adb"
	"github.com/sirupsen/logrus"
)

func StartServer() error {
	adb := adb.New()
	err := adb.Start()
	if err != nil {
		logrus.WithField("location", "StartServer").Error("failed start adb server: ", err)
	}
	router := NewRouter(adb)
	err = http.ListenAndServe(":3000", router.Setup())
	return err
}
