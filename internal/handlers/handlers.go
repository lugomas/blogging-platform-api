package handlers

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"roadmaps/projects/blogging-platform-api/internal/models"
)

func HandlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		models.GetAllPosts(w)
	case "POST":
		models.CreatePost(w, r)
	default:
		slog.Error("invalid method: ", "r.Method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]
	if postID == "" {
		slog.Error("invalid post id: ", "postID", postID)
		http.Error(w, "post id required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		models.GetPost(w, postID)
	case "PUT":
		models.UpdatePost(w, r, postID)
	case "DELETE":
		models.DeletePost(w, postID)
	default:
		slog.Error("invalid method: ", "r.Method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
