package routes

import (
	"net/http"

	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func MountAPIEndpoints(apiCfg *middleware.ApiConfig, router *chi.Mux) {
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		handlerReset(w, r, apiCfg)
	})

	apiRouter.Get("/chirps", func(w http.ResponseWriter, r *http.Request) {
		handlerGetChirps(w, r, apiCfg)
	})
	apiRouter.Post("/chirps", func(w http.ResponseWriter, r *http.Request) {
		handlerPostChirps(w, r, apiCfg)
	})
	apiRouter.Get("/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		handlerGetChirp(w, r, apiCfg)
	})

	apiRouter.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		handlerPostUsers(w, r, apiCfg)
	})
	apiRouter.Put("/users", func(w http.ResponseWriter, r *http.Request) {
		handlerPutUsers(w, r, apiCfg)
	})
	apiRouter.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		handlerLogin(w, r, apiCfg)
	})

	apiRouter.Post("/refresh", func(w http.ResponseWriter, r *http.Request) {
		handlerRefresh(w, r, apiCfg)
	})

	router.Mount("/api", apiRouter)
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func handlerReset(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	cfg.FileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
