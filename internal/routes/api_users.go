package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
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

	currentTime := time.Now().UTC()
	jwtCurrentTime := jwt.NewNumericDate(currentTime)

	var timeToExpire time.Time
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > 24*60*60 {
		timeToExpire = currentTime.Add(time.Hour * 24)
	} else {
		timeToExpire = currentTime.Add(time.Second * time.Duration(params.ExpiresInSeconds))
	}
	jwtTimeToExpire := jwt.NewNumericDate(timeToExpire)

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwtCurrentTime,
		ExpiresAt: jwtTimeToExpire,
		Subject:   strconv.Itoa(user.ID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(cfg.JWTSecret))
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
