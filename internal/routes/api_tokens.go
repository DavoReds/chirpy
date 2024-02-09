package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DavoReds/chirpy/internal/middleware"
)

func handlerRefresh(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	type response struct {
		Token string `json:"token"`
	}

	tokenString := extractAuthorizationHeader(r)
	if tokenString == "" {
		http.Error(w, "Missing refresh token", http.StatusBadRequest)
		return
	}

	wasRevoked, err := cfg.DB.WasTokenRevoked(tokenString)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if wasRevoked {
		http.Error(w, "Unathorized", http.StatusUnauthorized)
		return
	}

	token, err := parseJWT(tokenString, []byte(cfg.JWTSecret))
	if err != nil {
		http.Error(w, "Unathorized", http.StatusUnauthorized)
		return
	}

	isRefresh, err := verifyJWTIssuer(token, "chirpy-refresh")
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if !isRefresh {
		http.Error(w, "Not a refresh token", http.StatusUnauthorized)
		return
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	claims := getAccessJWTClaims(id)
	newToken, err := createJWT(claims, []byte(cfg.JWTSecret))
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})
}
