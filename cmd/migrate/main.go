package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

var (
	direction = flag.String("direction", "up", "migration direction: up, down, step, status")
	steps     = flag.Int("steps", 0, "number of steps to migrate when direction is step")
	force     = flag.Int("force", -1, "force a migration version and clear dirty state")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-direction up|down|step|status] [-steps N] [-force VERSION]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("database environment variables are not fully set")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}

	migrationsPath := filepath.Join(wd, "migrations")
	migrationsURL := (&url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(migrationsPath),
	}).String()

	m, err := migrate.New(
		migrationsURL,
		dsn,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if *force >= 0 {
		if err := m.Force(*force); err != nil {
			log.Fatalf("failed to force migration version: %v", err)
		}
		log.Printf("forced migration version to %d", *force)
		return
	}

	switch *direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "step":
		if *steps == 0 {
			log.Fatal("-steps must be provided when direction=step")
		}
		err = m.Steps(*steps)
	case "status":
		version, dirty, verr := m.Version()
		if verr != nil {
			log.Fatalf("failed to get migration version: %v", verr)
		}
		log.Printf("migration version: %d, dirty: %t", version, dirty)
		return
	default:
		usage()
		os.Exit(1)
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			version, dirty, verr := m.Version()
			if verr != nil {
				log.Fatalf("migration failed: %v", err)
			}
			log.Printf("no migration changes to apply (current version: %d, dirty: %t)", version, dirty)
			return
		}
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migration completed successfully")
}
