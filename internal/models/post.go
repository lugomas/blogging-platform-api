package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"roadmaps/projects/blogging-platform-api/internal/database"
	"time"
)

// Post represents the structure of a blog post
type Post struct {
	Id        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Category  string   `json:"category"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

// GetAllPosts retrieves all posts from the database and returns them as a JSON response
func GetAllPosts(w http.ResponseWriter) {
	slog.Info("GetAllPosts - fetching all posts")

	// Query all posts from the database
	rows, err := database.DB.Query("SELECT * FROM posts")
	if err != nil {
		slog.Error("Failed to fetch posts from database", "error", err)
		http.Error(w, "failed to fetch posts from database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post    // Slice to hold all retrieved posts
	var tagsJSON string // Temporary variable to hold tags in JSON format

	// Iterate through the query results
	for rows.Next() {
		var post Post
		// Scan each row into a Post struct
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Category, &tagsJSON, &post.CreatedAt, &post.UpdatedAt); err != nil {
			slog.Error("failed to parse posts: ", "error", err)
			http.Error(w, "failed to parse posts", http.StatusInternalServerError)
			return
		}

		// Parse the tags JSON string into the Tags slice
		err = json.Unmarshal([]byte(tagsJSON), &post.Tags)
		if err != nil {
			slog.Error("failed to unmarshal tags: ", "error", err)
			http.Error(w, "failed to unmarshal tags", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	// Send the retrieved posts as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

	slog.Info("GetAllPosts - Completed fetching posts from database")
}

// CreatePost creates a new post based on the request body and stores it in the database
func CreatePost(w http.ResponseWriter, r *http.Request) {
	slog.Info("CreatePost - Creating post")
	var post Post

	// Decode the request body into a Post struct
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		slog.Error("Request body decoding failed", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Convert the Tags slice into a JSON string
	tagsJSON, err := json.Marshal(post.Tags)
	if err != nil {
		slog.Error("failed to marshal tags: ", "error", err)
		http.Error(w, "failed to marshal tags", http.StatusInternalServerError)
		return
	}

	// Generate a UUID for the post
	post.Id = uuid.New().String()

	// Set timestamps for creation and last update
	post.CreatedAt = time.Now().Format("200601021504105")
	post.UpdatedAt = time.Now().Format("200601021504105")

	// Insert the new post into the database
	_, err = database.DB.Exec("INSERT INTO posts (id, title, content, category, tags, createdat, updatedat) VALUES (?, ?, ?, ?, ?, ?, ?)", post.Id, post.Title, post.Content, post.Category, tagsJSON, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		slog.Error("CreatePost - Failed to insert post into database", "error", err)
		http.Error(w, "failed to insert posts", http.StatusInternalServerError)
		return
	}

	// Respond with the created post
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)

	slog.Info("CreatePost - Post created successfully:", "id", post.Id)
}

// GetPost retrieves a specific post by ID and returns it as a JSON response
func GetPost(w http.ResponseWriter, id string) {
	slog.Info("GetPost - Fetching post with ID: ", "id", id)

	var post Post
	var tagsJSON string

	// Fetch the post from the database
	err := database.DB.QueryRow(
		"SELECT id, title, content, category, tags, createdat, updatedat FROM posts WHERE id = ?", id,
	).Scan(&post.Id, &post.Title, &post.Content, &post.Category, &tagsJSON, &post.CreatedAt, &post.UpdatedAt)

	// Handle no rows found.
	if errors.Is(err, sql.ErrNoRows) {
		slog.Error("post not found", "id", id)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Handle other SQL errors.
	if err != nil {
		slog.Error("failed to fetch posts: ", "error", err)
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	// Parse the tags JSON string into the Tags slice
	err = json.Unmarshal([]byte(tagsJSON), &post.Tags)
	if err != nil {
		slog.Error("failed to unmarshal tags: ", "error", err)
		http.Error(w, "failed to unmarshal tags", http.StatusInternalServerError)
		return
	}

	// Send the post as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)

	slog.Info("GetPost - Post fetched successfully:", "id", id)
}

// UpdatePost updates an existing post by ID with new data provided in the request body
func UpdatePost(w http.ResponseWriter, r *http.Request, id string) {
	slog.Info("UpdatePost - Updating post...", "id", id)

	var updatedPost Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		slog.Error("invalid request body: ", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tagsJSON, err := json.Marshal(updatedPost.Tags)
	if err != nil {
		slog.Error("failed to marshal tags: ", "error", err)
		http.Error(w, "failed to marshal tags", http.StatusInternalServerError)
		return
	}

	// Update post field
	updatedPost.UpdatedAt = time.Now().Format("200601021504105")
	result, err := database.DB.Exec("UPDATE posts SET id= ?, title = ?, content = ?, category = ?, tags = ?, updatedat = ? WHERE id = ?", id, updatedPost.Title, updatedPost.Content, updatedPost.Category, tagsJSON, updatedPost.UpdatedAt, id)
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

	// Send the post as a JSON response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPost)

	slog.Info("UpdatePost - Post updated", "id", id)
}

// DeletePost deletes a specific post by ID from the database
func DeletePost(w http.ResponseWriter, id string) {
	slog.Info("DeletePost - Deleting post...", "id", id)

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

	// Send the post as a JSON response
	w.WriteHeader(http.StatusNoContent)

	slog.Info("DeletePost - Post deleted", "id", id)
}
