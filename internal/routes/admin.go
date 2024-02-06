package routes

import (
	"fmt"
	"net/http"

	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func MountAdminEndpoints(apiCfg *middleware.ApiConfig, router *chi.Mux) {
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		HandlerMetrics(w, r, apiCfg)
	})

	router.Mount("/admin", adminRouter)
}

func HandlerMetrics(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
            <body>
                <h1>Welcome, Chirpy Admin</h1>
                <p>Chirpy has been visited %d times!</p>
            </body>
        </html>`, cfg.FileServerHits)))
}
