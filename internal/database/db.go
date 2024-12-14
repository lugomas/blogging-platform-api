package database

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

var DB *sql.DB

func DatabaseInit() {
	// Initialize the database connection without specifying a database
	var err error
	DB, err = sql.Open("mysql", "root:admin123@tcp(localhost:3306)/")
	if err != nil {
		slog.Error("Failed to connect to MySQL: ", "error", err)
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Test the database connection
	if err := DB.Ping(); err != nil {
		slog.Error("Failed to connect to MySQL: ", "error", err)
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Ensure the 'blogdb' database exists
	if err := createDatabase("blogdb"); err != nil {
		slog.Error("Failed to create database: ", "error", err)
		log.Fatalf("Failed to create database: %v", err)
	}

	// Reconnect to the 'blogdb' database
	DB.Close()
	DB, err = sql.Open("mysql", "root:admin123@tcp(localhost:3306)/blogdb")
	if err != nil {
		slog.Error("Failed to connect to blogdb: ", "error", err)
		log.Fatalf("Failed to connect to blogdb: %v", err)
	}

	// Test the database connection again
	if err := DB.Ping(); err != nil {
		slog.Error("Failed to connect to blogdb: ", "error", err)
		log.Fatalf("Failed to ping blogdb: %v", err)
	}

	// Create posts table if it doesn't exist
	if err := createPostsTable(); err != nil {
		slog.Error("Failed to create posts table: ", "error", err)
		log.Fatalf("Failed to create posts table: %v", err)
	}
}

// Create a database if it doesn't exist
func createDatabase(dbName string) error {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	_, err := DB.Exec(query)
	return err
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
