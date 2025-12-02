package repository

import (
	"context"

	"github.com/braunkc/todo-db/internal/domain/entities"
	valueobjects "github.com/braunkc/todo-db/internal/domain/value_objects/query"
)

type Repository interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
	DeleteUserByID(ctx context.Context, id string) error

	CreateTask(ctx context.Context, task *entities.Task) (*entities.Task, error)
	GetTaskByID(ctx context.Context, ID string) (*entities.Task, error)
	GetTasks(ctx context.Context, query *valueobjects.GetTasksQuery) ([]*entities.Task, int64, int64, error)
	UpdateTask(ctx context.Context, task *entities.Task) (*entities.Task, error)
	DeleteTasks(ctx context.Context, IDs []string) error
}
