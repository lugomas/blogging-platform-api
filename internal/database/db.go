package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

var (
	DB     *sql.DB
	dbName = "blogging05"
)

// DatabaseInit initializes the database connection and ensures the database and necessary tables exist
func DatabaseInit() {

	// Step 1: Establish an initial database connection without specifying a database
	slog.Info("initializing database ...")
	var err error
	DB, err = sql.Open("mysql", "root:admin123@tcp(localhost:3306)/")
	if err != nil {
		slog.Error("failed to connect to MySQL: ", "error", err)
	}

	// Step 2: Test the initial database connection
	slog.Info("verifying database connection...")
	if err := DB.Ping(); err != nil {
		slog.Error("failed to connect to MySQL: ", "error", err)
	}

	// Step 3: Ensure the specified database exists
	slog.Info("creating database")
	if err := createDatabase(dbName); err != nil {
		slog.Error("Failed to connect to database", "dbName", dbName, "error", err)
	}

	// Step 4: Reconnect to the database, now specifying the database name
	slog.Info("reconnecting to the database...")
	DB.Close()
	DB, err = sql.Open("mysql", "root:admin123@tcp(localhost:3306)/"+dbName)
	if err != nil {
		slog.Error("Failed to connect to database", "dbName", dbName, "error", err)
	}
	slog.Info("database reconnected")

	// Step 5: Verify the database connection again after reconnecting
	slog.Info("verifying database connection...")
	if err := DB.Ping(); err != nil {
		slog.Error("Failed to connect to database", "dbName", dbName, "error", err)
	}

	// Step 6: Ensure the required "posts" table exists
	slog.Info("creating table..")
	if err := createPostsTable(); err != nil {
		slog.Error("Failed to create posts table: ", "error", err)
	}
	slog.Info("table created")
}

// createDatabase creates the specified database if it does not already exist
func createDatabase(dbName string) error {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName)
	_, err := DB.Exec(query)
	return err
}

// createPostsTable creates the "posts" table if it does not already exist
func createPostsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id CHAR(36) PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
	    category VARCHAR(255) NOT NULL,
	    tags JSON NOT NULL,
	    createdat VARCHAR(255) NOT NULL,
	    updatedat VARCHAR(255) NOT NULL
	)`
	_, err := DB.Exec(query)
	return err
}
