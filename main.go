package main

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"roadmaps/projects/blogging-platform-api/internal/database"
	"roadmaps/projects/blogging-platform-api/internal/handlers"
)

func main() {

	// Initialize MySQL Client
	//dbAdress := os.Getenv("DATABASE_ADDR")
	database.DatabaseInit()

	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/posts", handlers.HandlePosts).Methods("GET", "POST")
	r.HandleFunc("/posts/{id}", handlers.HandlePost).Methods("GET", "PUT", "DELETE")

	// Start the HTTP server
	port := "8081"
	slog.Info("Server is starting", "port", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
