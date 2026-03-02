package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/revueexchange/api/internal/repository"
)

// Config holds all configuration
type Config struct {
	Environment string
	Port        int
	LogLevel    string

	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost string
	RedisPort int

	// JWT
	JWTSecret string

	// AWS
	AWSRegion    string
	AWSEndpoint  string
}

// Load loads configuration from environment
func Load() *Config {
	// Load .env file if exists
	_ = godotenv.Load()

	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnvInt("PORT", 8080),
		LogLevel:    getEnv("LOG_LEVEL", "info"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "revueadmin"),
		DBPassword: getEnv("DB_PASSWORD", "revueexchange"),
		DBName:     getEnv("DB_NAME", "revueexchange"),

		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnvInt("REDIS_PORT", 6379),

		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),

		AWSRegion:   getEnv("AWS_REGION", "us-east-1"),
		AWSEndpoint: getEnv("AWS_ENDPOINT", ""),
	}
}

// InitDB initializes database connection
func InitDB(cfg *Config) (*repository.Repository, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return repository.NewRepository(pool), nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		return intValue
	}
	return defaultValue
}
