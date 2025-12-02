package grpc

import (
	"context"

	"github.com/braunkc/todo-db/internal/application/usecases"
	pb "github.com/braunkc/todo-db/proto/database"
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
	GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error)
	UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error)
	DeleteTasksByID(ctx context.Context, req *pb.DeleteTasksByIDRequest) (*pb.DeleteTasksByIDResponse, error)
}

func New(usecasesService usecases.UsecasesService) GRPCServerService {
	return &grpcServerService{
		usecasesService: usecasesService,
	}
}

// TODO: write handlers

func (g *grpcServerService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

}

func (g *grpcServerService) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {

}

func (g *grpcServerService) DeleteUserByID(ctx context.Context, req *pb.DeleteUserByIDRequest) (*pb.DeleteUserByIDResponse, error) {

}

func (g *grpcServerService) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {

}

func (g *grpcServerService) GetTasks(ctx context.Context, req *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {

}

func (g *grpcServerService) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {

}

func (g *grpcServerService) DeleteTasksByID(ctx context.Context, req *pb.DeleteTasksByIDRequest) (*pb.DeleteTasksByIDResponse, error) {

}
