package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

var defaultCORSOrigins = []string{
	"http://localhost",
	"http://localhost:80",
	"http://localhost:5173",
	"http://127.0.0.1",
	"http://127.0.0.1:80",
	"http://127.0.0.1:5173",
}

type Config struct {
	DatabaseURL      string
	JWTSecret        string
	Port             string
	CORSOrigins      []string
	CORSAllowLocal   bool
	MigrationsPath   string
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

	origins, allowLocal := parseCORSOrigins(os.Getenv("CORS_ORIGINS"))

	return &Config{
		DatabaseURL:    dbURL,
		JWTSecret:      jwtSecret,
		Port:           port,
		CORSOrigins:    origins,
		CORSAllowLocal: allowLocal,
		MigrationsPath: migrationsPath,
	}, nil
}

func parseCORSOrigins(value string) ([]string, bool) {
	if strings.TrimSpace(value) == "" {
		return defaultCORSOrigins, true
	}

	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		if origin := strings.TrimSpace(part); origin != "" {
			origins = append(origins, origin)
		}
	}
	if len(origins) == 0 {
		return defaultCORSOrigins, true
	}
	return origins, false
}

// IsLocalOrigin matches browser origins used during local development.
func IsLocalOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	switch u.Hostname() {
	case "localhost", "127.0.0.1", "::1":
		return true
	default:
		return false
	}
}
