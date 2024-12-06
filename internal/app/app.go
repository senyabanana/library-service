package app

import (
	"database/sql"
	"net/http"

	"github.com/senyabanana/library-service/internal/config"
	"github.com/senyabanana/library-service/internal/handlers"
	"github.com/senyabanana/library-service/internal/logger"
	"github.com/senyabanana/library-service/internal/repository"
	"github.com/senyabanana/library-service/internal/router"
	"github.com/senyabanana/library-service/internal/services"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type App struct {
	Router *http.Server
	Logger *logger.Logger
}

func InitializeApp() (*App, error) {
	logg := logger.NewLogger()

	cfg, err := config.LoadConfig(".")
	if err != nil {
		logg.WithError(err).Fatal("Failed to load configuration")
	}

	logg.WithField("config", cfg).Info("Configuration loaded")

	runDBMigration(cfg.MigrationURL, cfg.DBConn, logg)

	db, err := sql.Open("postgres", cfg.DBConn)
	if err != nil {
		logg.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	repo := repository.NewSongRepository(db, logg)
	service := services.NewSongService(repo, logg)
	handler := handlers.NewSongHandler(service, logg)
	routes := router.SetupRoutes(handler, logg)

	return &App{
		Router: &http.Server{
			Addr:    ":8080",
			Handler: routes,
		},
		Logger: logg,
	}, nil
}

func runDBMigration(migrationURL, dBSource string, logg *logger.Logger) {
	migration, err := migrate.New(migrationURL, dBSource)
	if err != nil {
		logg.WithError(err).Fatal("Cannot create a new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		logg.WithError(err).Fatal("Failed to run migrate up")
	}

	logg.Info("Database migrated successfully")
}
