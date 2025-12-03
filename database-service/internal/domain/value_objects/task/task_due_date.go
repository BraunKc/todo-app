package valueobjects

import (
	"time"

	"github.com/braunkc/todo-app/database-service/pkg/errors"
)

type TaskDueDate int64

func NewDueDate(dueDate int64) (*TaskDueDate, error) {
	dd := TaskDueDate(dueDate)
	if err := dd.Validate(); err != nil {
		return nil, err
	}

	return &dd, nil
}

func (dd TaskDueDate) Validate() error {
	now := time.Now().UTC().Unix()
	due := int64(dd)
	if now >= due ||
		now+int64(time.Hour*24*30*12*100) <= due {
		return errors.ErrInvalidField
	}

	return nil
}
