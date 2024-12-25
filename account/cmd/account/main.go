package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"io/ioutil"

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

// initDbAccount initializes the database depending on the environment
func initDbAccount(cfg Config) {
	// If the environment is production, run the up.sql script to initialize the database
	if cfg.ENV == "prod" {
		log.Println("Running up.sql to initialize the database...")

		// Build the database connection string
		connStr := fmt.Sprintf("postgres://postgres:postgres@%s:5432/postgres", cfg.DatabaseURL)

		// Connect to the database
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Could not connect to the database: %v", err)
		}
		defer db.Close()

		// Read the SQL script (you would load your actual up.sql file here)
		sqlFile := "/usr/bin/up.sql" // Path to your SQL script file
		data, err := ioutil.ReadFile(sqlFile)
		if err != nil {
			log.Fatalf("Failed to read up.sql: %s", err)
		}

		// Execute the SQL script
		_, err = db.Exec(string(data))
		if err != nil {
			log.Fatalf("Failed to execute up.sql: %v", err,data)
		}

		log.Println("Database initialization complete.")
	}
}

func main() {
	// Load environment variables into the Config struct
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the database if needed
	initDbAccount(cfg)

	// Build the repository connection string
	repoConnStr := fmt.Sprintf("postgres://postgres:postgres@%s:5432/postgres", cfg.DatabaseURL)

	// Create a retryable repository connection
	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostgresRepository(repoConnStr)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	// Start the gRPC server
	log.Println("Listening on port 8080...")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}
