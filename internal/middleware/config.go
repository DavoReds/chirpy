package middleware

import (
	"net/http"
)

type ApiConfig struct {
	FileServerHits int
	Port           string
	FilesystemRoot string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits++
		next.ServeHTTP(w, r)
	})
}
