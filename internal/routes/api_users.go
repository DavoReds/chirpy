package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func handlerPostUsers(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	params := &parameters{}
	if err := decodeJSON(r.Body, params); err != nil {
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

	user, err := cfg.DB.CreateUser(params.Email, []byte(params.Password))
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
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	params := &parameters{}
	if err := decodeJSON(r.Body, params); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if params.Email == "" || params.Password == "" {
		log.Println("Request doesn't contain required fields")
		respondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	accessClaims := getAccessJWTClaims(user.ID)
	accessToken, err := createJWT(accessClaims, []byte(cfg.JWTSecret))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:    user.ID,
		Email: user.Email,
		Token: ss,
	})
}

func handlerPutUsers(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	tokenString := extractAuthorizationHeader(r)
	if tokenString == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Authorization header")
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	params := &parameters{}
	if err := decodeJSON(r.Body, params); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	user, err := cfg.DB.UpdateUser(id, params.Email, []byte(params.Password))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:    user.ID,
		Email: user.Email,
	})
}
