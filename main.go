package main

import (
	"log"
	"net/http"

	"github.com/DavoReds/chirpy/internal/database"
	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/DavoReds/chirpy/internal/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	apiCfg := middleware.ApiConfig{
		FileServerHits: 0,
		Port:           "8080",
		FilesystemRoot: ".",
		DB:             *database.NewDB("./database.json"),
	}

	r := chi.NewRouter()
	routes.MountAppEndpoints(&apiCfg, r)
	routes.MountAPIEndpoints(&apiCfg, r)
	routes.MountAdminEndpoints(&apiCfg, r)

	corsMux := middleware.MiddlewareCors(r)

	server := &http.Server{
		Addr:    ":" + apiCfg.Port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", apiCfg.FilesystemRoot, apiCfg.Port)
	log.Fatal(server.ListenAndServe())
}
