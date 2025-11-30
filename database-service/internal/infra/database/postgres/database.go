package database

import (
	"fmt"

	"gihtub.com/braunkc/todo-db/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type databaseService struct {
	db *gorm.DB
}

type DatabaseService interface {
}

func NewDatabaseService(cfg *config.Config) (DatabaseService, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to opeb database: %w", err)
	}

	return &databaseService{
		db: db,
	}, nil
}

// TODO: write CRUD funcs
