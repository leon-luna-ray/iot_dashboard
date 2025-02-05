package api

import (
	"net/http"
)

func Router() *http.ServeMux {
	router := http.NewServeMux()

	// API routes
	router.HandleFunc("/api/v1/hello", handleHello)
	router.HandleFunc("/api/v1/posts", handlePosts)
	router.HandleFunc("/api/v1/sensors", handleSensors)

	return router
}
