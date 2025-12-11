package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"prompt-library/database"
	"prompt-library/handlers"
	"prompt-library/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using defaults or environment variables")
	}

	fmt.Println("Starting Prompt Library API...")

	// initialize database
	database.InitDB()

	// create a router
	router := mux.NewRouter()

	// public routes, no auth needed
	router.HandleFunc("/api/login", handlers.Login).Methods("POST")
	router.HandleFunc("/api/prompts", handlers.GetAllPrompts).Methods("GET")
	router.HandleFunc("/api/audit-logs", handlers.GetAuditLogs).Methods("GET")

	// protected routes, needs auth
	// create sub router for authenticated routes

	router.Handle("/api/prompts", middleware.AuthMiddleware(http.HandlerFunc(handlers.CreatePrompt))).Methods("POST")
	router.Handle("/api/prompts/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.UpdatePrompt))).Methods("PUT")
	router.Handle("/api/prompts/{id}", middleware.AuthMiddleware(http.HandlerFunc(handlers.DeletePrompt))).Methods("DELETE")

	// enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	// wrap router with CORS
	handler := c.Handler(router)

	// Get port from .env or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
