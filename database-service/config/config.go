package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
	}
}

func New() (*Config, error) {
	var cfg Config

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
