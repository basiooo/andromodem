// Package templates provides functionality to embed and serve frontend assets
// (React JS build) directly from the Go application's filesystem.
package templates

import (
	"embed"
	"io/fs"
)

//go:embed andromodem_fe/dist/*
var MainPage embed.FS

func GetTemplateFS() fs.FS {
	return MainPage
}
