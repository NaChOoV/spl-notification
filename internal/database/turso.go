package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"spl-notification/internal/config"

	"github.com/pressly/goose/v3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"go.uber.org/fx"
)

func CreateTursoConnection(
	lc fx.Lifecycle,
	envConfig *config.EnvironmentConfig,
) *sql.DB {

	if envConfig.TursoAuthToken == "" || envConfig.TursoBaseUrl == "" {
		log.Fatal("TURSO_AUTH_TOKEN or TURSO_BASE_URL environment variables are not set")
	}

	dbURL := fmt.Sprintf("%s?authToken=%s", envConfig.TursoBaseUrl, envConfig.TursoAuthToken)
	db, err := sql.Open("libsql", dbURL)

	migrationErr := MigrateTurso(db)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			if err != nil {
				return fmt.Errorf("error opening connection to Turso: %w", err)
			}

			if migrationErr != nil {
				return fmt.Errorf("error running migrations: %w", migrationErr)
			}

			return nil
		},
		OnStop: func(context.Context) error {
			db.Close()
			return nil
		},
	})

	return db
}

func MigrateTurso(db *sql.DB) error {
	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// 4. Execute Migrations
	migrationsDir := "./migrations"
	log.Printf("Starting migrations from directory: %s", migrationsDir)

	// goose.Up applies all pending migrations.
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("error executing migrations with Goose: %w", err)
	}

	log.Printf("[MIGRATIONS] Applied successfully!")
	return nil
}
