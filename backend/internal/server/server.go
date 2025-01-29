package server

import (
	"log"
	"my-api/internal/api"
	"my-api/internal/config"
	"net/http"
)

func Start(cfg *config.Config) error {
	// Initialize the Homebridge client
	homebridgeClient := api.NewHomebridgeClient("http://192.168.50.2:8581", cfg.HBUsername, cfg.HBPassword)

	// Set up routes
	router := api.NewRouter(homebridgeClient)

	// Start the server
	log.Printf("Server started at http://localhost:%s 📡\n", cfg.Port)
	return http.ListenAndServe(":"+cfg.Port, router)
}
