package main

import (
	"fs/config"
	handlers "fs/server/internal"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	config.LoadConfig(env)

	// Ensure the uploaded_files directory exists
	ensureUploadDir()

	router := SetupRouter()

	log.Printf("Starting server on %s\n", config.ServerUrl)
	log.Fatal(http.ListenAndServe(config.ServerUrl, router))
}

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/upload", handlers.UploadFile)
	r.Get("/download/{fileName}", handlers.DownloadFile)

	return r
}

func ensureUploadDir() {
	const uploadDir = "./uploaded_files"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory: %s, error: %v", uploadDir, err)
		}
	}
}