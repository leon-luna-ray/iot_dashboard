package server

import (
	"iot_dashboard/internal/api"
	"iot_dashboard/internal/config"

	// "iot_dashboard/internal/homebridge"
	"log"
	"net/http"
)

func Start(cfg *config.Config) error {
	// homebridgeClient := homebridge.NewHomebridgeClient("http://localhost:8581", cfg.HBUsername, cfg.HBPassword)

	router := api.TestRouter()

	log.Printf("Server started at http://localhost:%s 📡\n", cfg.Port)
	return http.ListenAndServe(":"+cfg.Port, router)
}
