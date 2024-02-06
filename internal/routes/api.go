package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func MountAPIEndpoints(apiCfg *middleware.ApiConfig, router *chi.Mux) {
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Post("/chirps", handlerPostChirp)
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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func cleanString(text string) string {
	words := strings.Split(text, " ")

	for i, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle":
			words[i] = "****"
		case "sharbert":
			words[i] = "****"
		case "fornax":
			words[i] = "****"
		}
	}

	clean := strings.Join(words, " ")

	return clean
}

type validateChirpParams struct {
	Body string `json:"body"`
}

type validateChirpResponse struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := validateChirpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	body := cleanString(params.Body)
	respondWithJSON(w, 201, validateChirpResponse{Body: body})
}
