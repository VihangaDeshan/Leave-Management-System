package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"leave-management-system/internal/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg *config.Config) *sql.DB {
	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(5)                  // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Database connection established successfully")
	DB = db
	return db
}

// Close closes the database connection
func Close(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	if DB == nil {
		log.Fatal("Database not initialized. Call Connect() first.")
	}
	return DB
}
