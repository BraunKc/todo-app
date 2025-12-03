package valueobjects

import (
	"fmt"

	"github.com/braunkc/todo-app/database-service/pkg/errors"
	"github.com/google/uuid"
)

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

type SortField string

const (
	SortByPriority  SortField = "priority"
	SortByDueDate   SortField = "due_date"
	SortByCreatedAt SortField = "created_at"
)

type SortDirection string

const (
	SortAsc  SortDirection = "ASC"
	SortDesc SortDirection = "DESC"
)

type TaskFilters struct {
	Statuses   []TaskStatus
	Priorities []TaskPriority
}

type TaskOrderBy struct {
	Field     SortField
	Direction SortDirection
}

type GetTasksQuery struct {
	userID     string
	pageSize   int64
	pageNumber int64
	orderBy    TaskOrderBy
	filters    TaskFilters
	title      string
}

func NewGetTasksQuery(userID string, pageSize, pageNumber int64,
	sortField SortField, sortDirection SortDirection,
	statuses []TaskStatus, priorities []TaskPriority, title string) (*GetTasksQuery, error) {
	if pageSize < 1 || pageSize > 1000 {
		pageSize = 10
	}

	if pageNumber < 1 {
		pageNumber = 1
	}

	if sortField == "" {
		sortField = "priority"
	}

	if sortDirection == "" {
		sortDirection = "ASC"
	}

	if len(title) > 255 {
		title = title[:255]
	}

	query := GetTasksQuery{
		userID:     userID,
		pageSize:   pageSize,
		pageNumber: pageNumber,
		orderBy: TaskOrderBy{
			Field:     sortField,
			Direction: sortDirection,
		},
		filters: TaskFilters{
			Statuses:   statuses,
			Priorities: priorities,
		},
		title: title,
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	return &query, nil
}

func (q GetTasksQuery) Validate() error {
	if _, err := uuid.Parse(q.userID); err != nil {
		fmt.Println("uuid")
		return errors.ErrInvalidField
	}

	if !q.isValidSortField() {
		fmt.Println("field")
		return errors.ErrInvalidField
	}

	if !q.isValidSortDirection() {
		fmt.Println("dir")
		return errors.ErrInvalidField
	}

	return nil
}

func (q GetTasksQuery) isValidSortField() bool {
	if q.orderBy.Field == "" {
		fmt.Println("EMPTY FIELD")
	}

	switch q.orderBy.Field {
	case "", SortByPriority, SortByDueDate, SortByCreatedAt:
		return true
	default:
		return false
	}
}

func (q GetTasksQuery) isValidSortDirection() bool {
	switch q.orderBy.Direction {
	case SortAsc, SortDesc:
		return true
	default:
		return false
	}
}

func (q *GetTasksQuery) UserID() string {
	return q.userID
}

func (q *GetTasksQuery) PageSize() int64 {
	return q.pageSize
}

func (q *GetTasksQuery) PageNumber() int64 {
	return q.pageNumber
}

func (q *GetTasksQuery) Title() string {
	return q.title
}

func (q *GetTasksQuery) Filters() TaskFilters {
	return q.filters
}

func (q *GetTasksQuery) OrderBy() TaskOrderBy {
	return q.orderBy
}
