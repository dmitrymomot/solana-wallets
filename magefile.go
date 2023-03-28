//go:build mage
// +build mage

package main

import (
	"time"

	"github.com/fatih/color"
	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
	"github.com/magefile/mage/sh"
)

// Migrate database.
func Migrate() error {
	color.Cyan("Preparing migrations...")
	if err := sh.RunV("make", "prepare-migrations"); err != nil {
		return err
	}
	color.Yellow("Migrations copied to migrations directory.")

	color.Cyan("Migrating database...")
	return sh.RunV("go", "run", "./cmd/migrate/")
}

// Up runs the application in development mode (with database migrations).
func Up() error {
	color.Cyan("Starting application...")

	// Start database
	if err := sh.RunV("docker-compose", "up", "-d"); err != nil {
		return err
	}

	color.Yellow("Waiting for the database to be ready...")
	time.Sleep(5 * time.Second)

	// Migrate database
	if err := Migrate(); err != nil {
		return err
	}

	color.Cyan("Starting API server...")
	return sh.RunV("go", "run", "./cmd/api/")
}

// Down stops the application.
func Down() error {
	color.Cyan("Stopping application...")
	return sh.RunV("docker-compose", "down", "--rmi=local", "--volumes", "--remove-orphans")
}

// Run api
func Api() error {
	color.Cyan("Starting api...")
	return sh.RunV("go", "run", "./cmd/api/")
}

// Run tests
func Test() error {
	color.Cyan("Running tests...")
	return sh.RunV("go", "test", "-failfast", "-timeout 300s", "-p 1", "-count=1", "-race", "-cover", "-v", "./...")
}
