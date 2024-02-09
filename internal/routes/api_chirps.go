package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func handlerGetChirps(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	data, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	authorIDParam := r.URL.Query().Get("author_id")
	if authorIDParam != "" {
		id, err := strconv.Atoi(authorIDParam)
		if err != nil {
			http.Error(w, "Not a valid ID", http.StatusBadRequest)
			return
		}

		chirps, err := cfg.DB.GetChirpsFromAuthor(id)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}

func handlerPostChirps(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	type parameters struct {
		Body string `json:"body"`
	}

	id, err := getUserID(r, []byte(cfg.JWTSecret))
	if err != nil {
		log.Println(err)
		respondWithJSON(w, http.StatusUnauthorized, "Unathorized")
		return
	}

	params := &parameters{}
	if err := decodeJSON(r.Body, params); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	body := cleanString(params.Body)
	chirp, err := cfg.DB.CreateChirp(body, id)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func handlerGetChirp(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	idParam := chi.URLParam(r, "chirpID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.DB.GetChirpByID(id)
	if err != nil {
		http.Error(w, "Chirp doesn't exist", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func handlerDeleteChirp(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	idParam := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r, []byte(cfg.JWTSecret))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chirp, err := cfg.DB.GetChirpByID(chirpID)
	if err != nil {
		http.Error(w, "Chirp doesn't exist", http.StatusNotFound)
		return
	}

	if chirp.AuthorID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.DB.DeleteChirp(chirp.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
