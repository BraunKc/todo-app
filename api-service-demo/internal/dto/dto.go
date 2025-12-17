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
	ID string `json:""`
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
	ID          string `json:"id"`
	UserID      string
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	DueDate     int64        `json:"due_date"`
	CreatedAt   int64        `json:"created_at"`
}

type CreateTaskRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority"`
	DueDate     int64        `json:"due_date"`
}

type CreateTaskResponse struct {
	Task Task `json:"task"`
}

type GetTaskRequest struct {
	ID string `json:"id"`
}

type GetTaskResponse struct {
	Task Task
}

type Filters struct {
	TaskStatuses   []TaskStatus   `json:"task_statuses"`
	TaskPriorities []TaskPriority `json:"task_priorities"`
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
	Field     SortField     `json:"field"`
	Direction SortDirection `json:"direction"`
}

type GetTasksRequest struct {
	PageSize   int64   `json:"page_size"`
	PageNumber int64   `json:"page_number"`
	Filters    Filters `json:"filters"`
	OrderBy    OrderBy `json:"order_by"`
	Title      string  `json:"title"`
}

type GetTasksResponse struct {
	Tasks      []Task `json:"tasks"`
	TotalCount int64  `json:"total_count"`
	TotalPages int64  `json:"total_pages"`
}

type UpdateTaskRequest struct {
	ID          string        `json:"id"`
	Title       *string       `json:"title"`
	Description *string       `json:"description"`
	Status      *TaskStatus   `json:"status"`
	Priority    *TaskPriority `json:"priority"`
	DueDate     *int64        `json:"due_date"`
}

type UpdateTaskResponse struct {
	Task Task
}

type DeleteTasksByIDRequest struct {
	IDs []string `json:"ids"`
}

type DeleteTasksByIDResponse struct{}
