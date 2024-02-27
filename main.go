package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Could not load .env file %v\n", err)
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Empty Port String")
	}



	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept", "Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// v1Router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome! This is the homepage"))
	})
	// v1Router.Get("/errors", errorHandler)

	srv := http.Server{
		Addr:    ":" + portString,
		Handler: router,
	}
	fmt.Printf("This server runs on port: %v\n", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error: could not spin up server")
	}

}
