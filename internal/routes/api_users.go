package routes

import (
	"log"
	"net/http"
	"strconv"
	"strings"
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

func getClaims(expiresInSeconds int, userID int) *jwt.RegisteredClaims {
	currentTime := time.Now().UTC()
	jwtCurrentTime := jwt.NewNumericDate(currentTime)

	var timeToExpire time.Time
	if expiresInSeconds == 0 || expiresInSeconds > 24*60*60 {
		timeToExpire = currentTime.Add(time.Hour * 24)
	} else {
		timeToExpire = currentTime.Add(time.Second * time.Duration(expiresInSeconds))
	}
	jwtTimeToExpire := jwt.NewNumericDate(timeToExpire)

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwtCurrentTime,
		ExpiresAt: jwtTimeToExpire,
		Subject:   strconv.Itoa(userID),
	}

	return claims
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

	claims := getClaims(params.ExpiresInSeconds, user.ID)

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

func handlerPutUsers(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Authorization header")
		return
	}

	tokenString := strings.TrimPrefix(tokenHeader, "Bearer ")
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
