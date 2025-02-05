package server

import (
	"iot_dashboard/internal/config"
	"log"
	"net/http"
)

func Start(cfg *config.Config, router http.Handler) error {
	log.Printf("Server started at http://localhost:%s 📡\n", cfg.Port)
	log.Printf("QP_API_BASE: %s", cfg.QPAPIBase)
	return http.ListenAndServe(":"+cfg.Port, router)
}
