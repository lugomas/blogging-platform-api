package database

import (
	"database/sql"
	"log"
	"log/slog"

	_ "github.com/go-sql-driver/mysql" // Import MySQL drive
)

var DB *sql.DB

func DatabaseInit() {
	// Initialize the database connection
	var err error
	DB, err = sql.Open("mysql", "root:admin123@tcp(localhost:3306)/blogdb")
	if err != nil {
		slog.Error("Failed to connect to MySQL: ", "error", err)
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	//defer DB.Close()

	// Test the database connection
	if err := DB.Ping(); err != nil {
		slog.Error("Failed to connect to MySQL: ", "error", err)
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create posts table if it doesn't exist
	if err := createPostsTable(); err != nil {
		slog.Error("Failed to create posts table ", "error", err)
		log.Fatalf("Failed to create posts table: %v", err)
	}
}

// Create the posts table if it doesn't already exist
func createPostsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL
	)`
	_, err := DB.Exec(query)
	return err
}
