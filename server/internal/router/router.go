package router

import (
	"iot_dashboard/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.StripPrefix("/", handlers.StaticHandler()))
	return r
}
