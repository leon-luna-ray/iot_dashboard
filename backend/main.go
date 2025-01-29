package main

import (
	"embed"
	"iot_dashboard/internal/api"
	"iot_dashboard/internal/config"
	"iot_dashboard/internal/server"
	"log"
)

//go:embed public/static/dist/*
var staticFiles embed.FS

func main() {
	api.SetStaticFS(staticFiles)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Start the server
	if err := server.Start(cfg); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
