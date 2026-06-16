package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMigrationsPath = "migrations"
	maxConnectAttempts    = 30
	connectRetryDelay     = 2 * time.Second
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

// WaitForDatabase retries until Postgres accepts connections or the context is cancelled.
func WaitForDatabase(ctx context.Context, databaseURL string) error {
	var lastErr error
	for attempt := 1; attempt <= maxConnectAttempts; attempt++ {
		pool, err := Connect(ctx, databaseURL)
		if err == nil {
			pool.Close()
			if attempt > 1 {
				log.Printf("database ready after %d attempts", attempt)
			}
			return nil
		}
		lastErr = err
		log.Printf("waiting for database (attempt %d/%d): %v", attempt, maxConnectAttempts, err)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(connectRetryDelay):
		}
	}
	return fmt.Errorf("database not ready after %d attempts: %w", maxConnectAttempts, lastErr)
}

func RunMigrations(databaseURL, migrationsPath string) error {
	dir, err := resolveMigrationsPath(migrationsPath)
	if err != nil {
		return err
	}

	sourceURL := "file://" + filepath.ToSlash(dir)
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("migration version: %w", err)
	}
	if err == migrate.ErrNilVersion {
		log.Println("migrations: database schema is up to date (no migrations applied yet)")
	} else if dirty {
		return fmt.Errorf("migrations: database is in dirty state at version %d", version)
	} else {
		log.Printf("migrations: database schema is up to date (version %d)", version)
	}

	return nil
}

func resolveMigrationsPath(configured string) (string, error) {
	candidates := []string{}
	if configured != "" {
		candidates = append(candidates, configured)
	}
	candidates = append(candidates, defaultMigrationsPath)

	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), defaultMigrationsPath))
	}

	for _, candidate := range candidates {
		abs, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		if info, err := os.Stat(abs); err == nil && info.IsDir() {
			return abs, nil
		}
	}

	return "", fmt.Errorf("migrations directory not found (set MIGRATIONS_PATH)")
}
