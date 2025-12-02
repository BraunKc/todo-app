package usecases

import (
	"context"

	"github.com/braunkc/todo-db/internal/application/dto"
	"github.com/braunkc/todo-db/internal/application/repository"
	"github.com/braunkc/todo-db/internal/domain/entities"
	valueobjects "github.com/braunkc/todo-db/internal/domain/value_objects/query"
	"github.com/braunkc/todo-db/pkg/errors"
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

func (u *usecasesService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	user, err := entities.NewUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	resp, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.CreateUserResponse{
		User: dto.User{
			ID:           resp.ID(),
			Username:     resp.Username(),
			PasswordHash: string(resp.PasswordHash()),
		},
	}, nil
}

func (u *usecasesService) GetUserByUsername(ctx context.Context, req *dto.GetUserByUsernameRequest) (*dto.GetUserByUsernameResponse, error) {
	resp, err := u.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &dto.GetUserByUsernameResponse{
		User: dto.User{
			ID:           resp.ID(),
			Username:     resp.Username(),
			PasswordHash: string(resp.PasswordHash()),
		},
	}, nil
}

func (u *usecasesService) DeleteUserByID(ctx context.Context, req *dto.DeleteUserByIDRequest) (*dto.DeleteUserByIDResponse, error) {
	return &dto.DeleteUserByIDResponse{}, u.repo.DeleteUserByID(ctx, req.ID)
}

func (u *usecasesService) CreateTask(ctx context.Context, req *dto.CreateTaskRequest) (*dto.CreateTaskResponse, error) {
	userID := ctx.Value("userID")
	if userID == nil {
		return nil, errors.ErrFailedGetUserIDFromContext
	}

	task, err := entities.NewTask(userID.(string), req.Title, req.Description, 0, uint8(req.Priority), req.DueDate)
	if err != nil {
		return nil, err
	}

	resp, err := u.repo.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return &dto.CreateTaskResponse{
		Task: dto.Task{
			ID:          resp.ID(),
			Title:       resp.Title(),
			Description: resp.Description(),
			Status:      dto.TaskStatus(resp.Status()),
			Priority:    dto.TaskPriority(resp.Priority()),
			DueDate:     resp.DueDate(),
			CreatedAt:   resp.CreatedAt(),
		},
	}, nil
}

func (u *usecasesService) GetTasks(ctx context.Context, req *dto.GetTasksRequest) (*dto.GetTasksResponse, error) {
	userID := ctx.Value("userID")
	if userID == nil {
		return nil, errors.ErrFailedGetUserIDFromContext
	}

	var taskStatuses []valueobjects.TaskStatus
	for _, status := range req.Filters.TaskStatuses {
		if status <= dto.TaskStatus(valueobjects.TaskStatusDone) {
			taskStatuses = append(taskStatuses, valueobjects.TaskStatus(status))
		}
	}

	var taskPriorities []valueobjects.TaskPriority
	for _, priority := range req.Filters.TaskPriorities {
		if priority <= dto.TaskPriority(valueobjects.TaskPriorityHigh) {
			taskPriorities = append(taskPriorities, valueobjects.TaskPriority(priority))
		}
	}

	query, err := valueobjects.NewGetTasksQuery(userID.(string), req.PageSize, req.PageNumber,
		valueobjects.SortField(req.OrderBy.Field), valueobjects.SortDirection(req.OrderBy.Direction),
		taskStatuses, taskPriorities, req.Title)
	if err != nil {
		return nil, err
	}

	resp, totalCount, totalPages, err := u.repo.GetTasks(ctx, query)
	if err != nil {
		return nil, err
	}

	var tasks []dto.Task
	for _, task := range resp {
		tasks = append(tasks, dto.Task{
			ID:          task.ID(),
			Title:       task.Title(),
			Description: task.Description(),
			Status:      dto.TaskStatus(task.Status()),
			Priority:    dto.TaskPriority(task.Priority()),
			DueDate:     task.DueDate(),
			CreatedAt:   task.CreatedAt(),
		})
	}

	return &dto.GetTasksResponse{
		Tasks:      tasks,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (u *usecasesService) UpdateTask(ctx context.Context, req *dto.UpdateTaskRequest) (*dto.UpdateTaskResponse, error) {
	task, err := u.repo.GetTaskByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		if err := task.UpdateTitle(*req.Title); err != nil {
			return nil, err
		}
	}

	if req.Description != nil {
		if err := task.UpdateDescription(*req.Description); err != nil {
			return nil, err
		}
	}

	if req.Status != nil {
		if err := task.UpdateStatus(uint8(*req.Status)); err != nil {
			return nil, err
		}
	}

	if req.Priority != nil {
		if err := task.UpdatePriority(uint8(*req.Priority)); err != nil {
			return nil, err
		}
	}

	if req.DueDate != nil {
		if err := task.UpdateDueDate(*req.DueDate); err != nil {
			return nil, err
		}
	}

	task, err = u.repo.UpdateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateTaskResponse{
		Task: dto.Task{
			ID:          task.ID(),
			Title:       task.Title(),
			Description: task.Description(),
			Status:      dto.TaskStatus(task.Status()),
			Priority:    dto.TaskPriority(task.Priority()),
			DueDate:     task.DueDate(),
			CreatedAt:   task.CreatedAt(),
		},
	}, nil
}

func (u *usecasesService) DeleteTasks(ctx context.Context, req *dto.DeleteTasksByIDRequest) (*dto.DeleteTasksByIDResponse, error) {
	return &dto.DeleteTasksByIDResponse{}, u.repo.DeleteTasks(ctx, req.IDs)
}
