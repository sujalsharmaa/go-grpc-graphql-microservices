package main

import (
	"database/sql"
	"fmt"
	"log"
	"net" // Import net package for net.Listen
	"net/http"
	"time"
	"google.golang.org/grpc"
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

	repoConnStr := fmt.Sprintf("postgres://postgres:postgres@%s:5432/postgres", cfg.DatabaseURL)
	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostgresRepository(repoConnStr)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	// Start the gRPC server in a separate goroutine
	go func() {
		log.Println("Starting gRPC server on port 8080...")
		s := account.NewService(r)
		listener, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
		grpcServer := grpc.NewServer()
		account.RegisterServiceServer(grpcServer, s) // Make sure this function is available from the generated code
		log.Fatal(grpcServer.Serve(listener))
	}()

	// Simple "/" health check route on port 8080
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": 200, "message": "health ok"}`))
	})

	// Start HTTP server for health check on port 8080
	log.Println("Health check route available at / on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
