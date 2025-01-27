package handlers

import (
	"net/http"
	"path/filepath"
)

func StaticHandler() http.Handler {
	staticDir := filepath.Join("..", "web", "dist")
	return http.FileServer(http.Dir(staticDir))
}
