package middleware

import (
	"net/http"

	"github.com/DavoReds/chirpy/internal/database"
)

type ApiConfig struct {
	FileServerHits int
	Port           string
	FilesystemRoot string
	DB             database.DB
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits++
		next.ServeHTTP(w, r)
	})
}
