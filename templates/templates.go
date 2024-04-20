package templates

import "embed"

//go:embed andromodem-dashboard/dist/*
var mainPage embed.FS

func GetTemplateFS() embed.FS {
	return mainPage
}
