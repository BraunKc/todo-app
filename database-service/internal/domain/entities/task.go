package entities

import (
	"time"

	valueobjects "github.com/braunkc/todo-db/internal/domain/value_objects/task"
	"github.com/google/uuid"
)

type Task struct {
	id          string
	userID      string
	title       valueobjects.TaskTitle
	description valueobjects.TaskDescription
	status      valueobjects.TaskStatus
	priority    valueobjects.TaskPriority
	dueDate     valueobjects.TaskDueDate
	createdAt   int64
}

func NewTask(userID, title, description string,
	status, priority uint8, dueDate int64) (*Task, error) {
	t, err := valueobjects.NewTaskTitle(title)
	if err != nil {
		return nil, err
	}

	d, err := valueobjects.NewDescription(description)
	if err != nil {
		return nil, err
	}

	s, err := valueobjects.NewTaskStatus(status)
	if err != nil {
		return nil, err
	}

	p, err := valueobjects.NewTaskPriority(priority)
	if err != nil {
		return nil, err
	}

	dd, err := valueobjects.NewDueDate(dueDate)
	if err != nil {
		return nil, err
	}

	return &Task{
		id:          uuid.New().String(),
		userID:      userID,
		title:       *t,
		description: *d,
		status:      *s,
		priority:    *p,
		dueDate:     *dd,
		createdAt:   time.Now().Unix(),
	}, nil
}

func NewTaskFromStorage(id, userID, title, description string,
	status, priority uint8, dueDate, createdAt int64) *Task {
	return &Task{
		id:          id,
		userID:      userID,
		title:       valueobjects.TaskTitle(title),
		description: valueobjects.TaskDescription(description),
		status:      valueobjects.TaskStatus(status),
		priority:    valueobjects.TaskPriority(priority),
		dueDate:     valueobjects.TaskDueDate(dueDate),
		createdAt:   createdAt,
	}
}

func (t *Task) ID() string {
	return t.id
}

func (t *Task) UserID() string {
	return t.userID
}

func (t *Task) Title() string {
	return string(t.title)
}

func (t *Task) Description() string {
	return string(t.description)
}

func (t *Task) Status() uint8 {
	return uint8(t.status)
}

func (t *Task) Priority() uint8 {
	return uint8(t.priority)
}

func (t *Task) DueDate() int64 {
	return int64(t.dueDate)
}

func (t *Task) CreatedAt() int64 {
	return int64(t.createdAt)
}

func (t *Task) UpdateTitle(title string) error {
	newTitle, err := valueobjects.NewTaskTitle(title)
	if err != nil {
		return err
	}

	t.title = *newTitle

	return nil
}

func (t *Task) UpdateDescription(description string) error {
	newDescription, err := valueobjects.NewDescription(description)
	if err != nil {
		return err
	}

	t.description = *newDescription

	return nil
}

func (t *Task) UpdateStatus(status uint8) error {
	newStatus, err := valueobjects.NewTaskStatus(status)
	if err != nil {
		return err
	}

	t.status = *newStatus

	return nil
}

func (t *Task) UpdatePriority(priority uint8) error {
	newPriority, err := valueobjects.NewTaskPriority(priority)
	if err != nil {
		return err
	}

	t.priority = *newPriority

	return nil
}

func (t *Task) UpdateDueDate(dueDate int64) error {
	newDueDate, err := valueobjects.NewDueDate(dueDate)
	if err != nil {
		return err
	}

	t.dueDate = *newDueDate

	return nil
}
