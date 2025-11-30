package database

import (
	"github.com/braunkc/todo-db/internal/domain/entities"
	"github.com/braunkc/todo-db/internal/infra/database/postgres/models"
	"github.com/google/uuid"
)

type mapper struct{}

type Mapper interface {
	UserToModel(user *entities.User) (*models.User, error)
	UserToDomain(user *models.User) *entities.User
	TaskToModel(task *entities.Task) (*models.Task, error)
	TaskToDomain(task *models.Task) *entities.Task
}

func NewMapper() Mapper {
	return &mapper{}
}

func (r *mapper) UserToModel(user *entities.User) (*models.User, error) {
	id, err := uuid.Parse(user.ID())
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:           id,
		Username:     user.Username(),
		PasswordHash: user.PasswordHash(),
	}, nil
}

func (r *mapper) UserToDomain(user *models.User) *entities.User {
	return entities.NewUserFromStorage(user.ID.String(), user.Username, user.PasswordHash)
}

func (r *mapper) TaskToModel(task *entities.Task) (*models.Task, error) {
	id, err := uuid.Parse(task.ID())
	if err != nil {
		return nil, err
	}
	userID, err := uuid.Parse(task.UserID())
	if err != nil {
		return nil, err
	}

	return &models.Task{
		ID:          id,
		UserID:      userID,
		Title:       task.Title(),
		Description: task.Description(),
		Status:      task.Status(),
		Priority:    task.Priority(),
		DueDate:     task.DueDate(),
		CreatedAt:   task.CreatedAt(),
	}, nil
}

func (r *mapper) TaskToDomain(task *models.Task) *entities.Task {
	return entities.NewTaskFromStorage(task.ID.String(), task.UserID.String(),
		task.Title, task.Description, task.Status, task.Priority, task.DueDate)
}
