package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"vigilant/internal/migrations"
)

func runMigrationCommand() (handled bool, err error) {
	command := flag.String("cmd", "", "migration command: up or down (down rolls back one migration)")
	flag.Parse()

	if *command == "" {
		return false, nil
	}
	if *command != "up" && *command != "down" {
		return false, fmt.Errorf("unknown command %q; supported commands: up, down", *command)
	}

	databaseURL, err := getDatabaseURL()
	if err != nil {
		return true, err
	}

	migrationsDir := migrationDirectory()
	switch *command {
	case "up":
		err = migrations.Up(databaseURL, migrationsDir)
		if err == nil {
			log.Println("database migrations applied")
		}
	case "down":
		err = migrations.Down(databaseURL, migrationsDir)
		if err == nil {
			log.Println("one database migration rolled back")
		}
	}
	return true, err
}

func migrationDirectory() string {
	if dir := strings.TrimSpace(os.Getenv("MIGRATIONS_DIR")); dir != "" {
		return dir
	}
	return "migrations"
}
