package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:app/dist
var appDist embed.FS

func NewUIServer() http.Handler {
	appFS, _ := fs.Sub(appDist, "app/dist")
	return http.FileServer(http.FS(appFS))
}
