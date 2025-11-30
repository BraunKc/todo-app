package valueobjects

import "github.com/braunkc/todo-db/pkg/errors"

type TaskStatus uint8

const (
	TaskStatusTodo TaskStatus = iota
	TaskStatusInProgress
	TaskStatusDone
)

func NewTaskStatus(status uint8) (*TaskStatus, error) {
	s := TaskStatus(status)
	if ok := s.IsValid(); !ok {
		return nil, errors.ErrInvalidField
	}

	return &s, nil
}

func (s TaskStatus) IsValid() bool {
	return s <= TaskStatusDone
}
