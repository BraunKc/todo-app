package repository

import "github.com/braunkc/todo-db/internal/domain/entities"

type TaskRepository interface {
	CreateUser(user *entities.User) error
	GetUserByUsername(username string) (*entities.User, error)
	DeleteUserByID(id string) error

	CreateTask(task *entities.Task) error
	GetTaskByID(ID string) (*entities.Task, error)
	GetTasks(userID string) ([]*entities.Task, error)
	UpdateTask(task *entities.Task) error
	DeleteTasks(IDs []string) error
}
