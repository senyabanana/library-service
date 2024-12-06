package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/senyabanana/library-service/internal/config"
	"github.com/senyabanana/library-service/internal/handlers"
	"github.com/senyabanana/library-service/internal/repository"
	"github.com/senyabanana/library-service/internal/services"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	fmt.Println(cfg)

	runDBMigration(cfg.MigrationURL, cfg.DBConn)

	db, err := sql.Open("postgres", cfg.DBConn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := repository.NewSongRepository(db)
	service := services.NewSongService(repo)
	handler := handlers.NewSongHandler(*service)

	http.HandleFunc("/songs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetSongs(w, r)
		case http.MethodPost:
			handler.AddSong(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/songs/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			handler.DeleteSong(w, r)
		case http.MethodPut:
			handler.UpdateSong(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
