package main

import (
	"backend/internal/configs"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		printUsage()
		return
	}

	configs.LoadEnv()
	cfg := configs.Load()

	if cfg.DBHost == "" || cfg.DBName == "" {
		log.Fatal("Database configuration is missing. Ensure internal/configs/.env is loaded.")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		cfg.DBName,
		driver,
	)
	if err != nil {
		log.Fatalf("failed to initialize migrate instance: %v", err)
	}

	command := args[0]

	switch command {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("No migrations to apply (database is up to date).")
				return
			}
			log.Fatalf("migration up failed: %v", err)
		}
		fmt.Println("Successfully applied all pending migrations.")

	case "down":
		steps := 1
		if len(args) > 1 {
			parsedSteps, err := strconv.Atoi(args[1])
			if err != nil || parsedSteps < 1 {
				log.Fatalf("invalid steps argument: %s", args[1])
			}
			steps = parsedSteps
		}
		if err := m.Steps(-steps); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("No migrations to rollback.")
				return
			}
			log.Fatalf("migration down failed: %v", err)
		}
		fmt.Printf("Successfully rolled back %d migration(s).\n", steps)

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				fmt.Println("Current migration version: None (database has no migrations applied).")
				return
			}
			log.Fatalf("failed to get migration version: %v", err)
		}
		fmt.Printf("Current migration version: %d (dirty: %t)\n", version, dirty)

	case "force":
		if len(args) < 2 {
			log.Fatal("force command requires a target version argument, e.g.: go run ./cmd/migrate force 1")
		}
		v, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("invalid version argument: %s", args[1])
		}
		if err := m.Force(v); err != nil {
			log.Fatalf("failed to force version: %v", err)
		}
		fmt.Printf("Successfully forced migration version to %d.\n", v)

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Database Migration Tool (golang-migrate)")
	fmt.Println("\nUsage:")
	fmt.Println("  go run ./cmd/migrate up              # Apply all pending migrations")
	fmt.Println("  go run ./cmd/migrate down [steps]    # Rollback [steps] migrations (default 1)")
	fmt.Println("  go run ./cmd/migrate version         # Display current applied version")
	fmt.Println("  go run ./cmd/migrate force <version> # Force set migration version (e.g. to recover from dirty state)")
}
