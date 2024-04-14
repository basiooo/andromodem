package main

import (
	"github.com/basiooo/andromodem/internal/logger"
	"github.com/basiooo/andromodem/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.SetupLogger()
	err := server.StartServer()
	if err != nil {
		logrus.Fatalf("Failed start server: %v", err)
	}
}
