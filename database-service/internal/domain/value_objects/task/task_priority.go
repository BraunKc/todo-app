package valueobjects

import "github.com/braunkc/todo-app/database-service/pkg/errors"

type TaskPriority uint8

const (
	TaskPriorityLow TaskPriority = iota
	TaskPriorityMedium
	TaskPriorityHigh
)

func NewTaskPriority(priority uint8) (*TaskPriority, error) {
	p := TaskPriority(priority)
	if ok := p.IsValid(); !ok {
		return nil, errors.ErrInvalidField
	}

	return &p, nil
}

func (p TaskPriority) IsValid() bool {
	return p <= TaskPriorityHigh
}
