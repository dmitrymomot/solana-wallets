//go:build mage
// +build mage

package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
	"github.com/magefile/mage/sh"
)

// Prepare migrations
func PrepareMigrations() error {
	color.Cyan("Preparing migrations...")

	color.Yellow("Clean migrations directory...")
	if err := sh.RunV("rm", "-Rvf", "migrations/*.sql"); err != nil {
		return fmt.Errorf("failed to remove migrations: %w", err)
	}

	color.Yellow("Copy migrations...")
	if err := sh.RunV("cp", "-Rvf", "./svc/**/repository/sql/migrations/*.sql", "migrations/"); err != nil {
		color.Red("Failed to copy migrations! Error: %s", err.Error())
	}

	color.Yellow("Finished preparing migrations!")
	return nil
}

// Migrate database.
func Migrate() error {
	if err := PrepareMigrations(); err != nil {
		return err
	}

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
