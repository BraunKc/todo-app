package client

import (
	"context"
	"errors"

	"github.com/braunkc/todo-app/api-service-demo/internal/dto"
	pb "github.com/braunkc/todo-app/api-service-demo/proto/database"
	"golang.org/x/crypto/bcrypt"
)

type databaseService struct {
	client pb.DataBaseServiceClient
}

type DatabaseService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	Authenticate(ctx context.Context, username, password string) (*dto.User, error)
	DeleteUserByID(ctx context.Context, req *dto.DeleteUserByIDRequest) (*dto.DeleteUserByIDResponse, error)

	CreateTask(ctx context.Context, req *dto.CreateTaskRequest) (*dto.CreateTaskResponse, error)
	GetTask(ctx context.Context, req *dto.GetTaskRequest) (*dto.GetTaskResponse, error)
	GetTasks(ctx context.Context, req *dto.GetTasksRequest) (*dto.GetTasksResponse, error)
	UpdateTask(ctx context.Context, req *dto.UpdateTaskRequest) (*dto.UpdateTaskResponse, error)
	DeleteTasksByID(ctx context.Context, req *dto.DeleteTasksByIDRequest) (*dto.DeleteTasksByIDResponse, error)
}

func New(dbClient pb.DataBaseServiceClient) DatabaseService {
	return &databaseService{
		client: dbClient,
	}
}

func (db *databaseService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	resp, err := db.client.CreateUser(ctx, &pb.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &dto.CreateUserResponse{
		User: dto.User{
			ID:       resp.User.Id,
			Username: resp.User.Username,
		},
	}, nil
}

func (db *databaseService) Authenticate(ctx context.Context, username, password string) (*dto.User, error) {
	resp, err := db.client.GetUserByUsername(ctx, &pb.GetUserByUsernameRequest{
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(resp.User.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &dto.User{
		ID:       resp.User.Id,
		Username: resp.User.Username,
	}, nil
}

func (db *databaseService) DeleteUserByID(ctx context.Context, req *dto.DeleteUserByIDRequest) (*dto.DeleteUserByIDResponse, error) {
	_, err := db.client.DeleteUserByID(ctx, &pb.DeleteUserByIDRequest{
		Id: req.ID,
	})
	if err != nil {
		return nil, err
	}

	return &dto.DeleteUserByIDResponse{}, nil
}

func (db *databaseService) CreateTask(ctx context.Context, req *dto.CreateTaskRequest) (*dto.CreateTaskResponse, error) {
	resp, err := db.client.CreateTask(ctx, &pb.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		Priority:    pb.TaskPriority(req.Priority),
		DueDate:     req.DueDate,
	})
	if err != nil {
		return nil, err
	}

	return &dto.CreateTaskResponse{
		Task: mapTaskToDTO(resp.Task),
	}, nil
}

func (db *databaseService) GetTask(ctx context.Context, req *dto.GetTaskRequest) (*dto.GetTaskResponse, error) {
	resp, err := db.client.GetTask(ctx, &pb.GetTaskRequest{
		Id: req.ID,
	})
	if err != nil {
		return nil, err
	}

	return &dto.GetTaskResponse{
		Task: mapTaskToDTO(resp.Task),
	}, nil
}

func (db *databaseService) GetTasks(ctx context.Context, req *dto.GetTasksRequest) (*dto.GetTasksResponse, error) {
	taskStatuses := make([]pb.TaskStatus, 0, len(req.Filters.TaskStatuses))
	for _, status := range req.Filters.TaskStatuses {
		taskStatuses = append(taskStatuses, pb.TaskStatus(status))
	}

	taskPriorities := make([]pb.TaskPriority, 0, len(req.Filters.TaskPriorities))
	for _, priority := range req.Filters.TaskPriorities {
		taskPriorities = append(taskPriorities, pb.TaskPriority(priority))
	}

	resp, err := db.client.GetTasks(ctx, &pb.GetTasksRequest{
		PageSize:   req.PageSize,
		PageNumber: req.PageNumber,
		Filters: &pb.Filters{
			TaskStatuses:   taskStatuses,
			TaskPriorities: taskPriorities,
		},
		OrderBy: &pb.OrderBy{
			Field:     pb.SortField(req.OrderBy.Field),
			Direction: pb.SortDirection(req.OrderBy.Direction),
		},
		Title: &req.Title,
	})
	if err != nil {
		return nil, err
	}

	tasks := make([]dto.Task, 0, len(resp.Tasks))
	for _, task := range resp.Tasks {
		tasks = append(tasks, mapTaskToDTO(task))
	}

	return &dto.GetTasksResponse{
		Tasks:      tasks,
		TotalCount: resp.TotalCount,
		TotalPages: resp.TotalPages,
	}, nil
}

func ptr[T any](v T) *T {
	return &v
}

func (db *databaseService) UpdateTask(ctx context.Context, req *dto.UpdateTaskRequest) (*dto.UpdateTaskResponse, error) {
	var status *pb.TaskStatus
	if req.Status != nil {
		status = ptr(pb.TaskStatus(*req.Status))
	}

	var priority *pb.TaskPriority
	if req.Priority != nil {
		priority = ptr(pb.TaskPriority(*req.Priority))
	}

	resp, err := db.client.UpdateTask(ctx, &pb.UpdateTaskRequest{
		Id:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		DueDate:     req.DueDate,
	})
	if err != nil {
		return nil, err
	}

	return &dto.UpdateTaskResponse{
		Task: mapTaskToDTO(resp.Task),
	}, nil
}

func (db *databaseService) DeleteTasksByID(ctx context.Context, req *dto.DeleteTasksByIDRequest) (*dto.DeleteTasksByIDResponse, error) {
	_, err := db.client.DeleteTasksByID(ctx, &pb.DeleteTasksByIDRequest{
		Ids: req.IDs,
	})
	if err != nil {
		return nil, err
	}

	return &dto.DeleteTasksByIDResponse{}, nil
}

func mapTaskToDTO(t *pb.Task) dto.Task {
	return dto.Task{
		ID:          t.Id,
		UserID:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      dto.TaskStatus(t.Status),
		Priority:    dto.TaskPriority(t.Priority),
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
	}
}
