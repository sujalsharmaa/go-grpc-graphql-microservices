package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/akhilsharma90/go-graphql-microservice/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Config struct that reads environment variables
type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	ENV         string `envconfig:"ENV"`
}

// initDbAccount initializes the database for production environment
func initDbAccount(cfg Config) {
	if cfg.ENV == "prod" {
		log.Println("Initializing the database...")

		connStr := fmt.Sprintf("postgres://postgres:postgres@%s:5432/postgres", cfg.DatabaseURL)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Could not connect to the database: %v", err)
		}
		defer db.Close()

		// Execute SQL commands to create tables
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS accounts (
				id CHAR(27) PRIMARY KEY,
				name VARCHAR(24) NOT NULL
			);
		`)
		if err != nil {
			log.Fatalf("Failed to execute database initialization script: %v", err)
		}

		log.Println("Database initialization complete.")
	}
}

func main() {
	// Load configuration from environment variables
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database if necessary
	initDbAccount(cfg)

	var r account.Repository
	if cfg.ENV == "prod" {
		// Build the repository connection string
		repoConnStr := fmt.Sprintf("postgres://postgres:postgres@%s:5432/postgres", cfg.DatabaseURL)

		// Retry connecting to the repository
		retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
			r, err = account.NewPostgresRepository(repoConnStr)
			if err != nil {
				log.Println("Retrying connection to repository:", err)
			}
			return err
		})
		defer r.Close()
	}

	// Create and start the gRPC service
	s := account.NewService(r)
	log.Println("Listening on port 8080 for gRPC...")
	log.Fatal(account.ListenGRPC(s, 8080))
}
