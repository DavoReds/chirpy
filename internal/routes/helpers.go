package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
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
