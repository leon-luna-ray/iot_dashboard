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
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create the main router
	router := api.Router() // Get your API router

	// Add static files to the same router
	staticFS := getStaticFS()
	router.Handle("/", http.FileServer(http.FS(staticFS)))

	// Start the server with the unified router
	if err := server.Start(cfg, router); err != nil {
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
