package database

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

var (
	DB     *sql.DB
	dbName = "blogging03"
)

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

	// Ensure the database exists
	if err := createDatabase(dbName); err != nil {
		slog.Error("Failed to connect to database", "dbName", dbName, "error", err)
		log.Fatalf("Failed to create database: %v", err)
	}

	// Reconnect to the database
	DB.Close()
	DB, err = sql.Open("mysql", "root:admin123@tcp(localhost:3306)/"+dbName)
	if err != nil {
		slog.Error("Failed to connect to database", "dbName", dbName, "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the database connection again
	if err := DB.Ping(); err != nil {
		slog.Error("Failed to connect to database", "dbName", dbName, "error", err)
		log.Fatalf("Failed to ping database: %v", err)
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
		id CHAR(36) PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
	    category VARCHAR(255) NOT NULL,
	    tags JSON NOT NULL
	)`
	_, err := DB.Exec(query)
	return err
}
