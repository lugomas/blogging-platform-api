package main

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"roadmaps/projects/blogging-platform-api/internal/database"
	"roadmaps/projects/blogging-platform-api/internal/handlers"
)

func main() {

	// Initialize the database connection
	// This sets up the MySQL client used by the application
	database.DatabaseInit()

	// Create a new router using the Gorilla Mux library
	r := mux.NewRouter()

	// Define routes for the API
	// "/posts" endpoint supports:
	// - GET: Fetch all posts
	// - POST: Create a new post
	r.HandleFunc("/posts", handlers.HandlePosts).Methods("GET", "POST")

	// "/posts/{id}" endpoint supports:
	// - GET: Fetch a specific post by ID
	// - PUT: Update a specific post by ID
	// - DELETE: Delete a specific post by ID
	r.HandleFunc("/posts/{id}", handlers.HandlePost).Methods("GET", "PUT", "DELETE")

	// Start the HTTP server on the specified port
	port := "8081"
	slog.Info("Server is starting", "port", port)

	// Launch the server and handle errors if it fails to start
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("Server failed to start", "error", err)
	}

	// This log message confirms the server has started (unlikely to reach this point if an error occurs)
	slog.Info("Server has started", "port", port)
}
