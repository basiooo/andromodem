package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func SetupLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logFile, err := os.OpenFile("andromodem.log.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.WithField("location", "setupLogger").Fatal("error opening file: ", err)
	}
	logrus.SetOutput(logFile)
}
