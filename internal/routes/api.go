package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
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
	apiRouter.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		handlerLogin(w, r, apiCfg)
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

func handlerGetChirps(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	data, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, data)
	return
}

func handlerPostChirps(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	type params struct {
		Body string `json:"body" validate:"required"`
	}

	decoder := json.NewDecoder(r.Body)
	parameters := params{}
	if err := decoder.Decode(&parameters); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if len(parameters.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	body := cleanString(parameters.Body)
	chirp, err := cfg.DB.CreateChirp(body)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chirp, err := cfg.DB.GetChirp(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func handlerPostUsers(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	type parameters struct {
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required"`
	}

	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, password)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Email already used")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		ID:    user.ID,
		Email: user.Email,
	})
}

func handlerLogin(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	type parameters struct {
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required"`
	}

	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:    user.ID,
		Email: user.Email,
	})
}
