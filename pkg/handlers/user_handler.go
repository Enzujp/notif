package handlers

import (
	"net/http"

	"github.com/enzujp/notif/api"
	"github.com/enzujp/notif/database"
	"github.com/enzujp/notif/pkg/models"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// decode body
	if err := render.DecodeJSON(r.Body, &user); err != nil {
		render.Status(r, http.StatusBadRequest) // or 400 statusCode
		render.JSON(w, r, map[string]interface{}{"error": "Invalid request Payload"})
		return
	}

	// hash password before saving to database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{"error": "Failed to hash password"})
		return
	}
	if user.Password == "" {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]interface{}{"error": "Password field cannot be empty"})
	}
	user.Password = string(hashedPassword) // set user password to hashed password

	if err := database.DB.Create(&user).Error; err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{"error": "Failed to create user"})
		return
	}

	token, err := api.GenerateToken(&user)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]interface{}{"error": "Failed to generate token"})
		return
	}

	// rabbit MQ here
	render.JSON(w, r, map[string]interface{}{"token": token})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := render.DecodeJSON(r.Body, &user); err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]interface{}{"error": "Invalid request payload"})
		return
	}

	// check if user is an existing user
	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser); err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]interface{}{"error": "Invalid username or password"})
		
		return
	}
	//compare hashed passwords
	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]interface{}{"error": "Invalid username or password"})
		return
	}
	token, err := api.GenerateToken(&existingUser)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]interface{}{"error": "Failed to generate token."})
		return
	}
	render.JSON(w, r, map[string]interface{}{"token": token})
}
