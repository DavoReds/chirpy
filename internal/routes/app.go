package routes

import (
	"net/http"

	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func MountAppEndpoints(apiCfg *middleware.ApiConfig, router *chi.Mux) {
	fsHandler := apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(apiCfg.FilesystemRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)
}
