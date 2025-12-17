package grpc

import (
	"context"

	"github.com/braunkc/todo-app/database-service/internal/application/dto"
	"github.com/braunkc/todo-app/database-service/internal/application/usecases"
	pb "github.com/braunkc/todo-app/database-service/proto/database"
	"google.golang.org/grpc"
)

type grpcServerService struct {
	pb.UnimplementedDataBaseServiceServer
	usecasesService usecases.UsecasesService
}

type GRPCServerService interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
	GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error)
	DeleteUserByID(ctx context.Context, req *pb.DeleteUserByIDRequest) (*pb.DeleteUserByIDResponse, error)

	CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error)
	GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error)
	GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error)
	UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error)
	DeleteTasksByID(ctx context.Context, req *pb.DeleteTasksByIDRequest) (*pb.DeleteTasksByIDResponse, error)
}

func New(usecasesService usecases.UsecasesService) *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterDataBaseServiceServer(grpcServer, &grpcServerService{
		usecasesService: usecasesService,
	})

	return grpcServer
}

func (g *grpcServerService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	r := dto.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
	}

	resp, err := g.usecasesService.CreateUser(ctx, &r)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:           resp.User.ID,
			Username:     resp.User.Username,
			PasswordHash: resp.User.PasswordHash,
		},
	}, nil
}

func (g *grpcServerService) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	r := dto.GetUserByUsernameRequest{
		Username: req.Username,
	}

	resp, err := g.usecasesService.GetUserByUsername(ctx, &r)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserByUsernameResponse{
		User: &pb.User{
			Id:           resp.User.ID,
			Username:     resp.User.Username,
			PasswordHash: resp.User.PasswordHash,
		},
	}, nil
}

func (g *grpcServerService) DeleteUserByID(ctx context.Context, req *pb.DeleteUserByIDRequest) (*pb.DeleteUserByIDResponse, error) {
	r := dto.DeleteUserByIDRequest{
		ID: req.Id,
	}

	_, err := g.usecasesService.DeleteUserByID(ctx, &r)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteUserByIDResponse{}, nil
}

func (g *grpcServerService) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	r := dto.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		Priority:    dto.TaskPriority(req.Priority),
		DueDate:     req.DueDate,
	}

	resp, err := g.usecasesService.CreateTask(ctx, &r)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTaskResponse{
		Task: mapTaskToPB(resp.Task),
	}, nil
}

func (g *grpcServerService) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	r := dto.GetTaskRequest{
		ID: req.Id,
	}

	resp, err := g.usecasesService.GetTask(ctx, &r)
	if err != nil {
		return nil, err
	}

	return &pb.GetTaskResponse{
		Task: mapTaskToPB(resp.Task),
	}, nil
}

func (g *grpcServerService) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	taskStatuses := make([]dto.TaskStatus, 0)
	taskPriorities := make([]dto.TaskPriority, 0)
	if req.Filters != nil {
		taskStatuses = make([]dto.TaskStatus, 0, len(req.Filters.TaskStatuses))
		for _, status := range req.Filters.TaskStatuses {
			taskStatuses = append(taskStatuses, dto.TaskStatus(status))
		}

		taskPriorities = make([]dto.TaskPriority, 0, len(req.Filters.TaskPriorities))
		for _, priority := range req.Filters.TaskPriorities {
			taskPriorities = append(taskPriorities, dto.TaskPriority(priority))
		}
	}

	if req.Title == nil {
		title := ""
		req.Title = &title
	}

	if req.OrderBy == nil {
		orderBy := pb.OrderBy{}
		req.OrderBy = &orderBy
	}

	r := dto.GetTasksRequest{
		PageSize:   req.PageSize,
		PageNumber: req.PageNumber,
		Filters: dto.Filters{
			TaskStatuses:   taskStatuses,
			TaskPriorities: taskPriorities,
		},
		OrderBy: dto.OrderBy{
			Field:     dto.SortField(req.OrderBy.Field),
			Direction: dto.SortDirection(req.OrderBy.Direction),
		},
		Title: *req.Title,
	}

	resp, err := g.usecasesService.GetTasks(ctx, &r)
	if err != nil {
		return nil, err
	}

	var tasks []*pb.Task
	for _, task := range resp.Tasks {
		tasks = append(tasks, &pb.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      pb.TaskStatus(task.Status),
			Priority:    pb.TaskPriority(task.Priority),
			DueDate:     task.DueDate,
			CreatedAt:   task.CreatedAt,
		})
	}

	return &pb.GetTasksResponse{
		Tasks:      tasks,
		TotalCount: resp.TotalCount,
		TotalPages: resp.TotalPages,
	}, nil
}

func ptr[T any](v T) *T {
	return &v
}

func (g *grpcServerService) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	var status *dto.TaskStatus
	if req.Status != nil {
		switch *req.Status {
		case *pb.TaskStatus_TODO.Enum():
			status = ptr(dto.TaskStatusTodo)
		case *pb.TaskStatus_IN_PROGRESS.Enum():
			status = ptr(dto.TaskStatusInProgress)
		case *pb.TaskStatus_DONE.Enum():
			status = ptr(dto.TaskStatusDone)
		}
	}

	var priority *dto.TaskPriority
	if req.Priority != nil {
		switch *req.Priority {
		case *pb.TaskPriority_LOW.Enum():
			priority = ptr(dto.TaskPriorityLow)
		case *pb.TaskPriority_MEDIUM.Enum():
			priority = ptr(dto.TaskPriorityMedium)
		case *pb.TaskPriority_HIGH.Enum():
			priority = ptr(dto.TaskPriorityHigh)
		}
	}

	r := dto.UpdateTaskRequest{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		DueDate:     req.DueDate,
	}

	resp, err := g.usecasesService.UpdateTask(ctx, &r)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateTaskResponse{
		Task: mapTaskToPB(resp.Task),
	}, nil
}

func (g *grpcServerService) DeleteTasksByID(ctx context.Context, req *pb.DeleteTasksByIDRequest) (*pb.DeleteTasksByIDResponse, error) {
	r := dto.DeleteTasksByIDRequest{
		IDs: req.Ids,
	}

	_, err := g.usecasesService.DeleteTasks(ctx, &r)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteTasksByIDResponse{}, nil
}

func mapTaskToPB(t dto.Task) *pb.Task {
	return &pb.Task{
		Id:          t.ID,
		UserId:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Status:      pb.TaskStatus(t.Status),
		Priority:    pb.TaskPriority(t.Priority),
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
	}
}
