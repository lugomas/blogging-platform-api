package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

var (
	DB *sql.DB
	// Read environment variables
	user         = os.Getenv("MYSQL_USER")
	password     = os.Getenv("MYSQL_PASSWORD")
	rootPassword = os.Getenv("MYSQL_ROOT_PASSWORD")
	host         = os.Getenv("DB_HOST")
	port         = os.Getenv("DB_PORT")
	dbName       = os.Getenv("MYSQL_DATABASE")
	//dsn          = fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, password, host, port)
	rootDsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/", "root", rootPassword, host, port)
)

const maxRetries = 5
const retryInterval = 5 * time.Second

func connectToDatabase() (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		slog.Info("Attempting to connect to the database...", "rootDsn", rootDsn)
		db, err = sql.Open("mysql", rootDsn)
		if err != nil {
			slog.Error("Failed to open connection with database", "err", err)
			return nil, err
		}
		if err == nil {
			pingErr := db.Ping()
			if pingErr == nil {
				slog.Info("Successfully connected to the database!")
				return db, nil
			} else {
				slog.Warn("Failed to connect to the database, retrying...", "retry", i+1, "error", pingErr)
				return nil, pingErr
			}
		}
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("could not connect to the database after %d attempts: %w", maxRetries, err)
}

// DatabaseInit initializes the database connection and ensures the database and necessary tables exist
func DatabaseInit() {

	// Step 1: Establish an initial database connection without specifying a database
	slog.Info("Initializing database...")
	var err error
	DB, err = connectToDatabase()
	if err != nil {
		slog.Error("Exiting application due to database connection failure", "error", err)
		os.Exit(1)
	}

	// Step 2: Test the initial database connection
	slog.Info("verifying database connection...")
	if err := DB.Ping(); err != nil {
		slog.Error("failed to connect to MySQL: ", "error", err)
	}

	// Step 3: Ensure the specified database exists
	slog.Info("creating database")
	if err := createDatabase(dbName); err != nil {
		slog.Error("Failed to create database", "dbName", dbName, "error", err)
	}

	// Step 4: Reconnect to the database, now specifying the database name
	slog.Info("reconnecting to the database...")
	DB.Close()
	DB, err = sql.Open("mysql", rootDsn+dbName)
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
