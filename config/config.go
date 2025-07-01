package config

import (
	"log"
	"os"
)

type Config struct {
	SecretKey string
	DBConfig  DatabaseConfig
	Server    ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port    string
	Mode    string
	BaseURL string
}

func LoadConfig() *Config {
	var cfg Config

	cfg.SecretKey = os.Getenv("APP_SECRET_KEY")
	if cfg.SecretKey == "" {
		log.Fatalf("Error: Required environment variable APP_SECRET_KEY is not set. Application cannot start.")
	}

	cfg.DBConfig.Host = os.Getenv("DB_HOST")
	if cfg.DBConfig.Host == "" {
		log.Fatalf("Error: Required environment variable DB_HOST is not set. Application cannot start.")
	}

	cfg.DBConfig.Port = os.Getenv("DB_PORT")
	if cfg.DBConfig.Port == "" {
		log.Fatalf("Error: Required environment variable DB_PORT is not set. Application cannot start.")
	}

	cfg.DBConfig.User = os.Getenv("DB_USER")
	if cfg.DBConfig.User == "" {
		log.Fatalf("Error: Required environment variable DB_USER is not set. Application cannot start.")
	}

	cfg.DBConfig.Password = os.Getenv("DB_PASSWORD")
	if cfg.DBConfig.Password == "" {
		log.Fatalf("Error: Required environment variable DB_PASSWORD is not set. Application cannot start.")
	}

	cfg.DBConfig.Name = os.Getenv("DB_NAME")
	if cfg.DBConfig.Name == "" {
		log.Fatalf("Error: Required environment variable DB_NAME is not set. Application cannot start.")
	}

	cfg.DBConfig.SSLMode = os.Getenv("DB_SSL_MODE")
	if cfg.DBConfig.SSLMode == "" {
		log.Fatalf("Error: Required environment variable DB_SSL_MODE is not set. Application cannot start.")
	}

	cfg.Server.Port = os.Getenv("APP_SERVER_PORT")
	if cfg.Server.Port == "" {
		log.Fatalf("Error: Required environment variable APP_SERVER_PORT is not set. Application cannot start.")
	}

	cfg.Server.Mode = os.Getenv("APP_SERVER_MODE")
	if cfg.Server.Mode == "" {
		log.Fatalf("Error: Required environment variable APP_SERVER_MODE is not set. Application cannot start.")
	}

	cfg.Server.BaseURL = os.Getenv("APP_BASE_URL")
	if cfg.Server.BaseURL == "" {
		log.Fatalf("Error: Required environment variable APP_BASE_URL is not set. Application cannot start.")
	}

	log.Println("Configuration loaded successfully from environment variables.")
	return &cfg
}
