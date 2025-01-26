package router

import (
    "net/http"
    "iot_dashboard/internal/handlers"

    "github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
    r := mux.NewRouter()
    r.PathPrefix("/").Handler(http.StripPrefix("/", handlers.StaticHandler()))
    return r
}