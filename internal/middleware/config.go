package middleware

import (
	"fmt"
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

func (cfg *ApiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileServerHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *ApiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
            <body>
                <h1>Welcome, Chirpy Admin</h1>
                <p>Chirpy has been visited %d times!</p>
            </body>
        </html>`, cfg.FileServerHits)))
}
