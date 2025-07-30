package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Logging     LoggingConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	RootPass string
}

type LoggingConfig struct {
	Level string
}

func Load() (*Config, error) {
	env := getEnvWithDefault("ENV", "dev")
	
	config := &Config{
		Environment: env,
		Server: ServerConfig{
			Port: getEnvWithDefault("SERVER_PORT", "8080"),
			Host: getEnvWithDefault("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			User:     getEnvWithDefault("DB_USER", "taskuser"),
			Password: getEnvWithDefault("DB_PASS", "taskpass"),
			Host:     getEnvWithDefault("DB_HOST", "localhost"),
			Port:     getEnvWithDefault("DB_PORT", "3306"),
			Name:     getEnvWithDefault("DB", "taskapi"),
			RootPass: getEnvWithDefault("DB_ROOT_PASS", "rootpass"),
		},
		Logging: LoggingConfig{
			Level: getEnvWithDefault("LOG_LEVEL", "info"),
		},
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASS is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB is required")
	}

	if _, err := strconv.Atoi(c.Database.Port); err != nil {
		return fmt.Errorf("DB_PORT must be a valid integer: %w", err)
	}
	if _, err := strconv.Atoi(c.Server.Port); err != nil {
		return fmt.Errorf("SERVER_PORT must be a valid integer: %w", err)
	}

	return nil
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "dev"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "prod"
}

func (c *Config) ServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}