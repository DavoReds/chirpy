package routes

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/DavoReds/chirpy/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

func decodeJSON(reader io.Reader, readTo interface{}) error {
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(readTo); err != nil {
		return err
	}

	return nil
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

func extractBearerHeader(r *http.Request) string {
	tokenHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(tokenHeader, "Bearer ")

	return tokenString
}

func extractAPIKeyHeader(r *http.Request) string {
	tokenHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(tokenHeader, "ApiKey ")

	return tokenString
}

func getAccessJWTClaims(userID int) *jwt.RegisteredClaims {
	currentTime := time.Now().UTC()
	jwtCurrentTime := jwt.NewNumericDate(currentTime)

	timeToExpire := currentTime.Add(time.Hour)
	jwtTimeToExpire := jwt.NewNumericDate(timeToExpire)

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwtCurrentTime,
		ExpiresAt: jwtTimeToExpire,
		Subject:   strconv.Itoa(userID),
	}

	return claims
}

func getRefreshJWTClaims(userID int) *jwt.RegisteredClaims {
	currentTime := time.Now().UTC()
	jwtCurrentTime := jwt.NewNumericDate(currentTime)

	timeToExpire := currentTime.Add(time.Hour * 24 * 60)
	jwtTimeToExpire := jwt.NewNumericDate(timeToExpire)

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwtCurrentTime,
		ExpiresAt: jwtTimeToExpire,
		Subject:   strconv.Itoa(userID),
	}

	return claims
}

// Returns the string representation of a newly created JWT
func createJWT(claims *jwt.RegisteredClaims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func parseJWT(tokenString string, secret []byte) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func verifyJWTIssuer(token *jwt.Token, issuer string) (bool, error) {
	tokenIssuer, err := token.Claims.GetIssuer()
	if err != nil {
		return false, err
	}
	if tokenIssuer != issuer {
		return false, nil
	}

	return true, nil
}

func getUserID(r *http.Request, secret []byte) (int, error) {
	tokenString := extractBearerHeader(r)
	if tokenString == "" {
		return 0, errors.New("Authorization header not present")
	}

	token, err := parseJWT(tokenString, secret)
	if err != nil {
		return 0, err
	}

	isAccess, err := verifyJWTIssuer(token, "chirpy-access")
	if err != nil {
		return 0, err
	}
	if !isAccess {
		return 0, errors.New("Not an access token")
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func filterChirpsByAuthor(chirps []domain.Chirp, authorID int) []domain.Chirp {
	var chirpsByAuthor []domain.Chirp
	for _, chirp := range chirps {
		if chirp.AuthorID == authorID {
			chirpsByAuthor = append(chirpsByAuthor, chirp)
		}
	}

	return chirpsByAuthor
}

type sortOrder string

const (
	Ascending  sortOrder = "asc"
	Descending sortOrder = "desc"
)

func parseSortOrder(s string) (sortOrder, error) {
	switch strings.ToLower(s) {
	case "asc":
		return Ascending, nil
	case "desc":
		return Descending, nil
	default:
		return "", fmt.Errorf("Invalid order: %s", s)
	}
}

func sortChirps(chirps []domain.Chirp, order sortOrder) {
	slices.SortFunc(chirps, func(a, b domain.Chirp) int {
		if order == Ascending {
			return cmp.Compare(a.ID, b.ID)
		} else {
			return cmp.Compare(b.ID, a.ID)
		}
	})
}
