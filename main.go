package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/enzujp/notif/pkg/handlers"
	"github.com/enzujp/notif/rabbitmq"
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

	// middlewares
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome! This is the homepage"))
	})
	// routes
	v1Router.Post("/users/signup", handlers.Signup)
	v1Router.Post("/users/login", handlers.Login)

	router.Mount("/v1Router", v1Router)

	runRabbit()

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

func runRabbit() {
	if err := rabbitmq.InitRabbitMq(); err != nil {
		panic(fmt.Sprintf("Failed to intialize RabbitMq: %v", err))
	}
	defer rabbitmq.Close()
}
