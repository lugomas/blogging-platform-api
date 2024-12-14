package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"roadmaps/projects/blogging-platform-api/internal/database"
	"time"
)

type Post struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM posts")
	if err != nil {
		slog.Error("failed to fetch posts: ", "error", err)
		http.Error(w, "failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post); err != nil {
			slog.Error("failed to parse posts: ", "error", err)
			http.Error(w, "failed to parse posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		slog.Error("invalid request body: ", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	result, err := database.DB.Exec("INSERT INTO posts (id, title, content, category, tags) VALUES (?, ?, ?, ?)", post.Id, post.Title, post.Content, post.Category, post.Tags)
	if err != nil {
		slog.Error("failed to insert posts: ", "error", err)
		http.Error(w, "failed to insert posts", http.StatusInternalServerError)
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("failed to get last insert id for posts: ", "error", err)
		http.Error(w, "Failed to retrieve post ID", http.StatusInternalServerError)
		return
	}
	post.Id = fmt.Sprint(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func GetPost(w http.ResponseWriter, id string) {
	var post Post
	err := database.DB.QueryRow("SELECT * FROM posts WHERE id = ?", id).Scan(&post.Id, &post.Title, &post.Content)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Error("post not found", "id", id)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	} else if err != nil {
		slog.Error("failed to fetch posts: ", "error", err)
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request, id string) {
	var updatedPost Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		slog.Error("invalid request body: ", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("UPDATE posts SET id= ?, title = ?, content = ?, category = ?, tags = ? WHERE id = ?", updatedPost.Id, updatedPost.Title, updatedPost.Content, updatedPost.Category, updatedPost.Tags, id)
	if err != nil {
		slog.Error("failed to update posts: ", "error", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to get number of rows affected: ", "error", err)
		http.Error(w, "Failed to retrieve update status", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		slog.Error("failed to update post: no rows updated", "id", id)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPost)
}

func DeletePost(w http.ResponseWriter, id string) {
	result, err := database.DB.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		slog.Error("failed to delete posts: ", "error", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("failed to get number of rows affected: ", "error", err)
		http.Error(w, "Failed to retrieve delete status", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		slog.Error("failed to delete post: no rows deleted", "id", id)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
