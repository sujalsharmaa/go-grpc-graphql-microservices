package main

import (
	"log"
	"net/http" // Import the http package
	"time"
	"fmt"
	"github.com/akhilsharma90/go-graphql-microservice/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	ENV         string `envconfig:"ENV"`
}

// initDbAccount initializes the database depending on the environment
func initDbAccount(cfg Config) {
	if cfg.ENV == "prod" {
		log.Println("Running up.sql to initialize the database...")
		connStr := fmt.Sprintf("postgres://postgres:postgres@%s:5432/postgres", cfg.DatabaseURL)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Could not connect to the database: %v", err)
		}
		defer db.Close()

		db.Exec("CREATE TABLE IF NOT EXISTS accounts(id CHAR(27) PRIMARY KEY, name VARCHAR(24) NOT NULL);")
		if err != nil {
			log.Fatalf("Failed to execute up.sql: %v", err)
		}

		log.Println("Database initialization complete.")
	}
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	initDbAccount(cfg)

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	// Simple health check route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": 200, "message": "health ok"}`))
	})

	// Start HTTP server on port 8080 for health check
	go func() {
		log.Println("Starting HTTP server on port 8080 for health check...")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// gRPC server
	log.Println("Listening on port 8080 for gRPC...")
	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}
