package api

import (
	"io/fs"
	"log"
	"net/http"
)

var staticFS fs.FS

func SetStaticFS(embeddedFS fs.FS) {
	staticFS = embeddedFS
}

func Router() http.Handler {
	mux := http.NewServeMux()

	htmlContent, err := fs.Sub(staticFS, "public/static/dist")

	if err != nil {
		log.Panic(err)
	}
	fs := http.FileServer(http.FS(htmlContent))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "public/static/dist/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	}))

	// API routes
	mux.HandleFunc("/api/v1/hello", handleHello)
	mux.HandleFunc("/api/v1/posts", handlePosts)
	mux.HandleFunc("/api/v1/sensors", handleSensors)

	return mux
}
