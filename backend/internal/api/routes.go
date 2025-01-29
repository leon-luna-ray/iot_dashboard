package api

import (
	"embed"
	// "iot_dashboard/internal/homebridge"
	"net/http"
)

// go:embed internal/static/dist/*
var staticFiles embed.FS

func TestRouter() http.Handler {
	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.FS(staticFiles))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "iot_dashboard/internal/static/dist/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	}))

	// API routes
	mux.HandleFunc("/api/v1/hello", handleHello)
	mux.HandleFunc("/api/v1/posts", handlePosts)
	// mux.HandleFunc("/api/v1/homebridge", func(w http.ResponseWriter, r *http.Request) {
	// 	homebridge.handleHomebridge(w, r, homebridgeClient)
	// })

	return mux
}
