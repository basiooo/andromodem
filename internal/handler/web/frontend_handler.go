package web

import (
	"embed"
	"io/fs"
	"net/http"
	"text/template"

	"github.com/basiooo/andromodem/templates"
	"go.uber.org/zap"
)

type FrontendHandler struct {
	Logger     *zap.Logger
	FrontendFS embed.FS
}

func NewFrontendHandler(logger *zap.Logger, frontendFS embed.FS) IFrontendHandler {
	return &FrontendHandler{
		Logger:     logger,
		FrontendFS: frontendFS,
	}
}

func (f *FrontendHandler) ServeAssets() http.Handler {
	staticFs, err := fs.Sub(templates.MainPage, "andromodem_fe/dist/assets")
	if err != nil {
		f.Logger.Fatal("failed to prepare static assets fs", zap.Error(err))
	}
	return http.StripPrefix("/assets/", http.FileServer(http.FS(staticFs)))
}

func (f *FrontendHandler) ServeIndex(writer http.ResponseWriter, request *http.Request) {
	frontEndTemplate, err := template.ParseFS(templates.MainPage, "andromodem_fe/dist/index.html")
	if err != nil {
		f.Logger.Error("error parsing template", zap.Error(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = frontEndTemplate.Execute(writer, nil)
	if err != nil {
		f.Logger.Error("error executing template", zap.Error(err))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
