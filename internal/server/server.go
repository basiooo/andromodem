package server

import (
	"net/http"

	"github.com/basiooo/andromodem/internal/adb"
	"github.com/basiooo/andromodem/templates"
	"github.com/sirupsen/logrus"
)

func StartServer() error {
	adb := adb.New()
	err := adb.Start()
	if err != nil {
		logrus.WithField("location", "StartServer").Error("failed start adb server: ", err)
	}
	mainPage := templates.GetTemplateFS()
	router := NewRouter(adb, mainPage)
	err = http.ListenAndServe(":49153", router.Setup())
	return err
}
