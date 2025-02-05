package server

import (
	"iot_dashboard/internal/config"
	"log"
	"net"
	"net/http"
)

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "localhost"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func Start(cfg *config.Config, router http.Handler) error {
	ip := getOutboundIP()
	log.Printf("Server started 📡")
	log.Printf("(💻 Local):  http://localhost:%s", cfg.Port)
	log.Printf("(🌐 Network):  http://%s:%s", ip, cfg.Port)
	return http.ListenAndServe(":"+cfg.Port, router)
}
