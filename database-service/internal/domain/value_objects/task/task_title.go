package valueobjects

import (
	"strings"

	"github.com/braunkc/todo-db/pkg/errors"
)

type TaskTitle string

func NewTaskTitle(title string) (*TaskTitle, error) {
	t := TaskTitle(title)
	if err := t.Validate(); err != nil {
		return nil, err
	}

	return &t, nil
}

func (t TaskTitle) Validate() error {
	if strings.TrimSpace(string(t)) == "" {
		return errors.ErrEmptyField
	}

	if len(t) > 128 {
		return errors.ErrTooLongField
	}

	return nil
}
