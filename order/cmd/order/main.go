package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/akhilsharma90/go-graphql-microservice/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Config struct that reads environment variables
type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL"`
	ENV         string `envconfig:"ENV"`
}

// initDbOrder initializes the database for the production environment
func initDbOrder(cfg Config) {
	if cfg.ENV == "prod" {
		log.Println("Initializing the database...")

		connStr := fmt.Sprintf("postgres://postgres:postgres@%s:5432/postgres", cfg.DatabaseURL)

		var db *sql.DB
		var err error

		// Retry connecting to the database
		retry.ForeverSleep(2*time.Second, func(_ int) error {
			db, err = sql.Open("postgres", connStr)
			if err != nil {
				log.Println("Retrying database connection:", err)
				return err
			}

			// Ping the database to ensure the connection is valid
			if pingErr := db.Ping(); pingErr != nil {
				log.Println("Retrying database ping:", pingErr)
				return pingErr
			}
			return nil
		})
		defer db.Close()

		// Execute SQL script to create tables
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS orders (
				id CHAR(27) PRIMARY KEY,
				created_at TIMESTAMP WITH TIME ZONE NOT NULL,
				account_id CHAR(27) NOT NULL,
				total_price MONEY NOT NULL
			);
			CREATE TABLE IF NOT EXISTS order_products (
				order_id CHAR(27) REFERENCES orders(id) ON DELETE CASCADE,
				product_id CHAR(27),
				quantity INT NOT NULL,
				PRIMARY KEY (product_id, order_id)
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
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the database if necessary
	initDbOrder(cfg)

	// Build the repository connection string
	repoConnStr := fmt.Sprintf(cfg.DatabaseURL)

	// Retry connecting to the repository
	var r order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewPostgresRepository(repoConnStr)
		if err != nil {
			log.Println("Retrying connection to repository:", err)
		}
		return err
	})
	defer r.Close()

	// Create and start the gRPC service
	s := order.NewService(r)
	log.Println("Listening on port 8080 for gRPC...")
	log.Fatal(order.ListenGRPC(s, cfg.AccountURL, cfg.CatalogURL, 8080))
}
