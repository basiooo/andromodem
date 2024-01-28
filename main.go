package main

import (
	"embed"
	adb "github.com/abccyz/goadb"
	"github.com/basiooo/andromodem/app"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/exec"
)

//go:embed andromodem-frontend/dist/*
var templateFS embed.FS

func setupADBClient() *adb.Adb {
	adbClient, err := adb.New()
	if err != nil {
		logrus.Fatalf("failed to create ADB client: %v", err)
	}

	err = adbClient.StartServer()
	if err != nil {
		logrus.Info("trying to start a ADB server manualy")
		cmd := exec.Command("adb", "start-server")
		err = cmd.Run()
		if err != nil {
			logrus.Fatalf("failed to start ADB server: %v", err)
		}
		logrus.Info("successfully run ADB server")
	}
	return adbClient
}

func setupLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	f, err := os.OpenFile("andromodemlog.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.WithField("function", "main").Fatal("error opening file: ", err)
	}
	logrus.SetOutput(f)
}

func main() {
	setupLogger()
	adbClient := setupADBClient()
	r := app.NewRouter(templateFS, adbClient)
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		logrus.WithField("function", "main").Fatal("failed start server : ", err)
	}
}
