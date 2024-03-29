package routes

import (
	"log"
	"net/http"
	"os"

	"github.com/DavoReds/chirpy/internal/middleware"
)

func handlerPolkaWebhook(w http.ResponseWriter, r *http.Request, cfg *middleware.ApiConfig) {
	type parameters struct {
		Data struct {
			UserID int `json:"user_id"`
		} `json:"data"`
		Event string `json:"event"`
	}

	polkaKey := os.Getenv("POLKA_KEY")
	apiKey := extractAPIKeyHeader(r)
	if apiKey == "" || apiKey != polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := &parameters{}
	if err := decodeJSON(r.Body, params); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if params.Event != "user.upgraded" {
		return
	}

	if err := cfg.DB.UpgradeUser(params.Data.UserID); err != nil {
		if err.Error() == "User doesn't exist" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	return
}
