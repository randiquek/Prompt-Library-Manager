package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key to sign tokens
// Should be environment variable in prod

// Returns jwt secret from .env
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable not set!")
	}
	return []byte(secret)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// get credentials from .env
	validUsername := os.Getenv("ADMIN_USERNAME")
	validPassword := os.Getenv("ADMIN_PASSWORD")

	if validUsername == "" || validPassword == "" {
		log.Println("Warning: admin username or password is not set in environment")
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	// check if entered login matches stored credentials
	if loginReq.Username != validUsername || loginReq.Password != validPassword {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	//create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginReq.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), //expires in 24hrs

	})

	// sign token with secret

	tokenString, err := token.SignedString(GetJWTSecret())
	if err != nil {
		log.Println("Error signing token", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// return token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}
