package main

import (
	"github.com/basiooo/andromodem/internal/server"
	"github.com/sirupsen/logrus"
)

func main() {
	err := server.StartServer()
	if err != nil {
		logrus.Fatalf("Failed start server: %v", err)
	}
}
