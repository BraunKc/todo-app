package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	HTTPServer struct {
		Port string `yaml:"port"`
	} `yaml:"http-server"`
	DatabaseService struct {
		GRPCAddr string
	}
	SecretKey string
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

	cfg.DatabaseService.GRPCAddr = os.Getenv("GRPC_ADDR")
	cfg.SecretKey = os.Getenv("SECRET_KEY")

	return &cfg, nil
}
