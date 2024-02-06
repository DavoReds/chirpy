package routes

import (
	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func MountAPIEndpoints(apiCfg *middleware.ApiConfig, router *chi.Mux) {
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", HandlerReadiness)
	apiRouter.HandleFunc("/reset", apiCfg.HandleReset)

	router.Mount("/api", apiRouter)
}
