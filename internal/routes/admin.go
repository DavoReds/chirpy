package routes

import (
	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func MountAdminEndpoints(apiCfg *middleware.ApiConfig, router *chi.Mux) {
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.HandlerMetrics)

	router.Mount("/admin", adminRouter)
}
