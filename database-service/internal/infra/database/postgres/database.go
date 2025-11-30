package database

import (
	"fmt"

	"github.com/braunkc/todo-db/config"
	"github.com/braunkc/todo-db/internal/application/repository"
	"github.com/braunkc/todo-db/internal/domain/entities"
	"github.com/braunkc/todo-db/internal/infra/database/postgres/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type databaseRepository struct {
	db     *gorm.DB
	mapper Mapper
}

func NewDatabaseService(cfg *config.Config, mapper Mapper) (repository.TaskRepository, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to opeb database: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate user: %w", err)
	}
	if err := db.AutoMigrate(&models.Task{}); err != nil {
		return nil, fmt.Errorf("failed to migrate task: %w", err)
	}

	return &databaseRepository{
		db:     db,
		mapper: mapper,
	}, nil
}

func (r *databaseRepository) CreateUser(user *entities.User) error {
	u, err := r.mapper.UserToModel(user)
	if err != nil {
		return err
	}

	return r.db.Create(u).Error
}

func (r *databaseRepository) GetUserByUsername(username string) (*entities.User, error) {
	var user models.User
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return r.mapper.UserToDomain(&user), nil
}

func (r *databaseRepository) DeleteUserByID(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id).Delete(&models.Task{}).Error; err != nil {
			return err
		}

		return tx.Where("id = ?", id).Delete(&models.User{}).Error
	})
}

func (r *databaseRepository) CreateTask(task *entities.Task) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id = ?", task.UserID()).First(&models.User{}).Error; err != nil {
			return err
		}

		t, err := r.mapper.TaskToModel(task)
		if err != nil {
			return err
		}

		return tx.Create(t).Error
	})
}

func (r *databaseRepository) GetTaskByID(ID string) (*entities.Task, error) {
	var t models.Task
	if err := r.db.Where("id = ?", ID).First(&t).Error; err != nil {
		return nil, err
	}

	return r.mapper.TaskToDomain(&t), nil
}

func (r *databaseRepository) GetTasks(userID string) ([]*entities.Task, error) {
	var t []models.Task
	if err := r.db.Where("user_id = ?", userID).Find(&t).Error; err != nil {
		return nil, err
	}

	var tasks []*entities.Task
	for _, task := range t {
		tasks = append(tasks, r.mapper.TaskToDomain(&task))
	}

	return tasks, nil
}

func (r *databaseRepository) UpdateTask(task *entities.Task) error {
	t, err := r.mapper.TaskToModel(task)
	if err != nil {
		return err
	}

	return r.db.Save(t).Error
}

func (r *databaseRepository) DeleteTasks(IDs []string) error {
	return r.db.Where("id IN ?", IDs).Delete(&models.Task{}).Error
}
