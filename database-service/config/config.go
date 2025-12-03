package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
)

type Config struct {
	GRPCServer struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc-server"`
	Database struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
	}
}

func New() (*Config, error) {
	file, err := os.ReadFile("./config/config.yml")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.Port = os.Getenv("DB_PORT")
	cfg.Database.Name = os.Getenv("DB_NAME")
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")

	return &cfg, nil
}
