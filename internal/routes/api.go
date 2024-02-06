package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func MountAPIEndpoints(apiCfg *middleware.ApiConfig, router *chi.Mux) {
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Post("/chirps", func(w http.ResponseWriter, r *http.Request) {
		handlerPostChirp(w, r, apiCfg)
	})
	apiRouter.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		HandleReset(w, r, apiCfg)
	})

	router.Mount("/api", apiRouter)
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func HandleReset(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	cfg.FileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func handlerPostChirp(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	type params struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	parameters := params{}
	if err := decoder.Decode(&parameters); err != nil {
		log.Println(err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(parameters.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	body := cleanString(parameters.Body)
	chirp, err := cfg.DB.CreateChirp(body)
	if err != nil {
		log.Println(err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	respondWithJSON(w, 201, chirp)
}
