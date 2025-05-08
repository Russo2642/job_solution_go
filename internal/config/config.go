package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	PostgreSQL PostgreSQLConfig
	JWT        JWTConfig
	Security   SecurityConfig
	RateLimit  RateLimitConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type PostgreSQLConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Secret           string
	ExpiresIn        time.Duration
	RefreshExpiresIn time.Duration
}

type SecurityConfig struct {
	PasswordSalt string
}

type RateLimitConfig struct {
	Requests int
	Duration time.Duration
}

func Load() (*Config, error) {
	godotenv.Load()

	serverPort := getEnv("SERVER_PORT", "8080")
	serverMode := getEnv("SERVER_MODE", "debug")

	pgHost := getEnv("POSTGRES_HOST", "localhost")
	pgPort := getEnv("POSTGRES_PORT", "5432")
	pgUser := getEnv("POSTGRES_USER", "postgres")
	pgPassword := getEnv("POSTGRES_PASSWORD", "postgres")
	pgDatabase := getEnv("POSTGRES_DB", "jobsolution")
	pgSSLMode := getEnv("POSTGRES_SSLMODE", "disable")

	pgMaxOpenConns, err := strconv.Atoi(getEnv("POSTGRES_MAX_OPEN_CONNS", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_MAX_OPEN_CONNS: %w", err)
	}

	pgMaxIdleConns, err := strconv.Atoi(getEnv("POSTGRES_MAX_IDLE_CONNS", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_MAX_IDLE_CONNS: %w", err)
	}

	pgConnMaxLifetime, err := time.ParseDuration(getEnv("POSTGRES_CONN_MAX_LIFETIME", "5m"))
	if err != nil {
		return nil, fmt.Errorf("invalid POSTGRES_CONN_MAX_LIFETIME: %w", err)
	}

	jwtSecret := getEnv("JWT_SECRET", "default_jwt_secret")
	jwtExpiresIn, err := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRES_IN: %w", err)
	}
	jwtRefreshExpiresIn, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRES_IN", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRES_IN: %w", err)
	}

	passwordSalt := getEnv("PASSWORD_SALT", "default_password_salt")

	rateLimitRequests, err := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_REQUESTS: %w", err)
	}
	rateLimitDuration, err := time.ParseDuration(getEnv("RATE_LIMIT_DURATION", "1m"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_DURATION: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: serverPort,
			Mode: serverMode,
		},
		PostgreSQL: PostgreSQLConfig{
			Host:            pgHost,
			Port:            pgPort,
			User:            pgUser,
			Password:        pgPassword,
			Database:        pgDatabase,
			SSLMode:         pgSSLMode,
			MaxOpenConns:    pgMaxOpenConns,
			MaxIdleConns:    pgMaxIdleConns,
			ConnMaxLifetime: pgConnMaxLifetime,
		},
		JWT: JWTConfig{
			Secret:           jwtSecret,
			ExpiresIn:        jwtExpiresIn,
			RefreshExpiresIn: jwtRefreshExpiresIn,
		},
		Security: SecurityConfig{
			PasswordSalt: passwordSalt,
		},
		RateLimit: RateLimitConfig{
			Requests: rateLimitRequests,
			Duration: rateLimitDuration,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
