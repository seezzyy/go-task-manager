package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL     string
	HTTPPort        string
	ShutdownTimeout time.Duration
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func Load() (Config, error) {
	databaseURL := getEnv("DATABASE_URL", "")

	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	httpPort := getEnv("HTTP_PORT", "8080")

	port, err := strconv.Atoi(httpPort)

	if err != nil {
		return Config{}, fmt.Errorf("некорректный HTTP_PORT: %w", err)
	}

	if port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("HTTP_PORT должен быть между 1 и 65535")
	}

	timeoutStr := getEnv("SHUTDOWN_TIMEOUT", "5s")

	duration, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return Config{}, fmt.Errorf("Некорректный ShutdownTimeout: %w", err)
	}

	if duration <= 0 {
		return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT должен быть больше нуля")
	}

	return Config{
		DatabaseURL:     databaseURL,
		HTTPPort:        httpPort,
		ShutdownTimeout: duration,
	}, nil
}
