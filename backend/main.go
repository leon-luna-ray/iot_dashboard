package main

import (
	"iot_dashboard/internal/config"
	"iot_dashboard/internal/server"
	"log"
)

func main() {
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
