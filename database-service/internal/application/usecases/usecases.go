package usecases

import (
	"context"

	"github.com/braunkc/todo-db/internal/application/dto"
	"github.com/braunkc/todo-db/internal/application/repository"
)

type usecasesService struct {
	repo repository.Repository
}

type UsecasesService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	GetUserByUsername(ctx context.Context, req *dto.GetUserByUsernameRequest) (*dto.GetUserByUsernameResponse, error)
	DeleteUserByID(ctx context.Context, req *dto.DeleteUserByIDRequest) (*dto.DeleteUserByIDResponse, error)

	CreateTask(ctx context.Context, req *dto.CreateTaskRequest) (*dto.CreateTaskResponse, error)
	GetTasks(ctx context.Context, req *dto.GetTasksRequest) (*dto.GetTasksResponse, error)
	UpdateTask(ctx context.Context, req *dto.UpdateTaskRequest) (*dto.UpdateTaskResponse, error)
	DeleteTasks(ctx context.Context, req *dto.DeleteTasksByIDRequest) (*dto.DeleteTasksByIDResponse, error)
}

func NewUsecasesService(repo repository.Repository) UsecasesService {
	return &usecasesService{
		repo: repo,
	}
}

// TODO: write CRUD funcs

func (u *usecasesService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {

}

func (u *usecasesService) GetUserByUsername(ctx context.Context, req *dto.GetUserByUsernameRequest) (*dto.GetUserByUsernameResponse, error) {

}

func (u *usecasesService) DeleteUserByID(ctx context.Context, req *dto.DeleteUserByIDRequest) (*dto.DeleteUserByIDResponse, error) {

}

func (u *usecasesService) CreateTask(ctx context.Context, req *dto.CreateTaskRequest) (*dto.CreateTaskResponse, error) {

}

func (u *usecasesService) GetTasks(ctx context.Context, req *dto.GetTasksRequest) (*dto.GetTasksResponse, error) {

}

func (u *usecasesService) UpdateTask(ctx context.Context, req *dto.UpdateTaskRequest) (*dto.UpdateTaskResponse, error) {

}

func (u *usecasesService) DeleteTasks(ctx context.Context, req *dto.DeleteTasksByIDRequest) (*dto.DeleteTasksByIDResponse, error) {

}
