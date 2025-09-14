package main

import (
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN is not set in environment")
	}
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("could not init migrate: %v", err)
	}
	if len(os.Args) < 2 {
		log.Fatal("You must provide a command: up, down, drop, version, step")
	}
	switch os.Args[1] {
	case "up":
		err = m.Up()
	case "down":
		err = m.Steps(-1)
	case "drop":
		err = m.Drop()
	case "version":
		var v uint
		v, _, err = m.Version()
		if err == nil {
			log.Printf("Current version: %d\n", v)
		}
	case "step":
		if len(os.Args) < 3 {
			log.Fatal("You must provide a number for step")
		}
		step, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid step value: %v", err)
		}
		err = m.Steps(step)
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
	if err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
	log.Println("Migration executed successfully ✅")
}
