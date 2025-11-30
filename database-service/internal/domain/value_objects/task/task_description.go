package valueobjects

import (
	"strings"

	"github.com/braunkc/todo-db/pkg/errors"
)

type TaskDescription string

func NewDescription(description string) (*TaskDescription, error) {
	d := TaskDescription(strings.TrimSpace(description))
	if err := d.Validate(); err != nil {
		return nil, err
	}

	return &d, nil
}

func (d TaskDescription) Validate() error {
	description := string(d)
	if len(description) > 1024 {
		return errors.ErrTooLongField
	}

	return nil
}
