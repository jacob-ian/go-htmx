package htmx

import (
	"embed"
	"net/http"
)

//go:embed assets
var staticFs embed.FS

func NewStaticFileServer() http.Handler {
	return http.FileServer(http.FS(staticFs))
}
