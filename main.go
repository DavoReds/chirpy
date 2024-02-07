package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/DavoReds/chirpy/internal/database"
	"github.com/DavoReds/chirpy/internal/middleware"
	"github.com/DavoReds/chirpy/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	debug_mode := flag.Bool("debug", false, "Run the program in debug mode")
	flag.Parse()

	if *debug_mode {
		os.Remove("./database.json")
	}

	jwtSecret, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		log.Fatal("JWT_SECRET env variable not set")
	}

	apiCfg := middleware.ApiConfig{
		FileServerHits: 0,
		Port:           "8080",
		FilesystemRoot: ".",
		DB:             *database.NewDB("./database.json"),
		JWTSecret:      jwtSecret,
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
