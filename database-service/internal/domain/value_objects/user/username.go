package valueobjects

import (
	"fmt"
	"strings"

	"github.com/braunkc/todo-db/pkg/errors"
)

type Username string

func NewUsername(username string) (*Username, error) {
	u := Username(username)
	if err := u.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate username: %w", err)
	}

	return &u, nil
}

func (u Username) Validate() error {
	if strings.TrimSpace(string(u)) == "" {
		return errors.ErrEmptyField
	}

	if len(u) > 64 {
		return errors.ErrTooLongField
	}

	return nil
}
