package main

import (
	"embed"
	"io/fs"
	"iot_dashboard/internal/api"
	"iot_dashboard/internal/config"
	"iot_dashboard/internal/server"
	"log"
	"net/http"
)

//go:embed public/dist/*
var staticFiles embed.FS

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create the router
	router := api.Router()

	// Add static files
	staticFS := getStaticFS()
	router.Handle("/", http.FileServer(http.FS(staticFS)))

	// CORS
	corsHandler := enableCors(router)

	// Start the server
	if err := server.Start(cfg, corsHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getStaticFS() fs.FS {
	subFS, err := fs.Sub(staticFiles, "public/dist")
	if err != nil {
		log.Fatal("Failed to create sub filesystem: ", err)
	}
	return subFS
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for frontend development
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
