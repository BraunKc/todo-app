package dto

type User struct {
	ID           string
	Username     string
	PasswordHash string
}

type CreateUserRequest struct {
	Username string
	Password string
}

type CreateUserResponse struct {
	User User
}

type GetUserByUsernameRequest struct {
	Username string
}

type GetUserByUsernameResponse struct {
	User User
}

type DeleteUserByIDRequest struct {
	ID string
}

type DeleteUserByIDResponse struct{}

type TaskStatus uint8

const (
	TaskStatusTodo TaskStatus = iota
	TaskStatusInProgress
	TaskStatusDone
)

type TaskPriority uint8

const (
	TaskPriorityLow TaskPriority = iota
	TaskPriorityMedium
	TaskPriorityHigh
)

type Task struct {
	ID          string
	Title       string
	Description string
	Status      TaskStatus
	Priority    TaskPriority
	DueDate     int64
	CreatedAt   int64
}

type CreateTaskRequest struct {
	Title       string
	Description string
	Priority    TaskPriority
	DueDate     int64
}

type CreateTaskResponse struct {
	Task Task
}

type Filters struct {
	TaskStatuses   []TaskStatus
	TaskPriorities []TaskPriority
}

type SortField uint8

const (
	Priority SortField = iota
	DueDate
	CreatedAt
)

type SortDirection uint8

const (
	Asc SortDirection = iota
	Desc
)

type OrderBy struct {
	Field     SortField
	Direction SortDirection
}

type GetTasksRequest struct {
	PageSize   int64
	PageNumber int64
	Filters    Filters
	OrderBy    OrderBy
	Title      string
}

type GetTasksResponse struct {
	Tasks      []Task
	TotalCount int64
	TotalPages int64
}

type UpdateTaskRequest struct {
	ID          string
	Title       *string
	Description *string
	Status      *TaskStatus
	Priority    *TaskPriority
	DueDate     *int64
}

type UpdateTaskResponse struct {
	Task Task
}

type DeleteTasksByIDRequest struct {
	IDs []string
}

type DeleteTasksByIDResponse struct{}
