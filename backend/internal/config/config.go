package config

import (
	"fmt"
	"os"
	"strings"
)

var defaultCORSOrigins = []string{
	"http://localhost:5173",
	"http://localhost:80",
	"http://localhost",
}

type Config struct {
	DatabaseURL    string
	JWTSecret      string
	Port           string
	CORSOrigins    []string
	MigrationsPath string
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}

	return &Config{
		DatabaseURL:    dbURL,
		JWTSecret:      jwtSecret,
		Port:           port,
		CORSOrigins:    parseCORSOrigins(os.Getenv("CORS_ORIGINS")),
		MigrationsPath: migrationsPath,
	}, nil
}

func parseCORSOrigins(value string) []string {
	if strings.TrimSpace(value) == "" {
		return defaultCORSOrigins
	}

	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		if origin := strings.TrimSpace(part); origin != "" {
			origins = append(origins, origin)
		}
	}
	if len(origins) == 0 {
		return defaultCORSOrigins
	}
	return origins
}
