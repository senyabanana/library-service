package main

import (
	"fmt"
	"log"

	"library-service/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	fmt.Println(cfg)

	runDBMigration(cfg.MigrationURL, cfg.DBConn)
}

func runDBMigration(migrationURL, dBSource string) {
	migration, err := migrate.New(migrationURL, dBSource)
	if err != nil {
		log.Fatal("cannot create a new migrate instance", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up:", err)
	}
	log.Println("db migrated successfully")
}
