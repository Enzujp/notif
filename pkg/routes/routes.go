package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/enzujp/notif/pkg/handlers"
)

func UserRoutes(router chi.Router) {
	router.Post("/user/signup", handlers.Signup)
	router.Post("/user/login", handlers.Login)
}