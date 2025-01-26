// filepath: /Users/rayluna/code/projects/iot_dashboard/server/cmd/main.go
package main

import (
    "log"
    "net/http"
    "iot_dashboard/internal/router"
)

func main() {
    r := router.SetupRouter()
    log.Println("Server is running on http://localhost:9090")
    log.Fatal(http.ListenAndServe(":9090", r))
}