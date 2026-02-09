package config

import (
	"os"
	"strconv"

	pc "github.com/nightmaker00/accordion-go/pkg/db/postgres"
)

type Config struct {
	Server struct {
		Address  string
		Port     string
		Timeouts struct {
			ReadSeconds  int
			WriteSeconds int
			IdleSeconds  int
		}
	}
	pc.Config
}

func Load() (*Config, error) {
	cfg := &Config{}

	//default address
	cfg.Server.Address = "0.0.0.0"
	cfg.Server.Port = "8080"
	cfg.Server.Timeouts.ReadSeconds = 5
	cfg.Server.Timeouts.WriteSeconds = 10
	cfg.Server.Timeouts.IdleSeconds = 60

	cfg.Config.Host = "localhost"
	cfg.Config.Port = "5432"
	cfg.Config.User = "postgres"
	cfg.Config.Password = "postgres"
	cfg.Config.DBName = "faq"
	cfg.Config.SSLMode = "disable"

	//in env if exists
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.Server.Address = host
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}
	if seconds, ok := getEnvInt("SERVER_READ_TIMEOUT_SECONDS"); ok {
		cfg.Server.Timeouts.ReadSeconds = seconds
	}
	if seconds, ok := getEnvInt("SERVER_WRITE_TIMEOUT_SECONDS"); ok {
		cfg.Server.Timeouts.WriteSeconds = seconds
	}
	if seconds, ok := getEnvInt("SERVER_IDLE_TIMEOUT_SECONDS"); ok {
		cfg.Server.Timeouts.IdleSeconds = seconds
	}

	if host := os.Getenv("POSTGRES_HOST"); host != "" {
		cfg.Config.Host = host
	}
	if port := os.Getenv("POSTGRES_PORT"); port != "" {
		cfg.Config.Port = port
	}
	if user := os.Getenv("POSTGRES_USER"); user != "" {
		cfg.Config.User = user
	}
	if password := os.Getenv("POSTGRES_PASSWORD"); password != "" {
		cfg.Config.Password = password
	}
	if name := os.Getenv("POSTGRES_DB"); name != "" {
		cfg.Config.DBName = name
	}
	if ssl := os.Getenv("POSTGRES_SSLMODE"); ssl != "" {
		cfg.Config.SSLMode = ssl
	}

	return cfg, nil
}

func getEnvInt(key string) (int, bool) {
	raw := os.Getenv(key)
	if raw == "" {
		return 0, false
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	if value < 0 {
		return 0, false
	}
	return value, true
}
