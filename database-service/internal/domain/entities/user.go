package entities

import (
	"fmt"

	valueobjects "github.com/braunkc/todo-db/internal/domain/value_objects/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id           string
	username     valueobjects.Username
	passwordHash []byte
}

func NewUser(username, password string) (*User, error) {
	u, err := valueobjects.NewUsername(username)
	if err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	return &User{
		id:           uuid.New().String(),
		username:     *u,
		passwordHash: passwordHash,
	}, nil
}

func NewUserFromStorage(id, username string, passwordHash []byte) *User {
	return &User{
		id:           id,
		username:     valueobjects.Username(username),
		passwordHash: passwordHash,
	}
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Username() string {
	return string(u.username)
}

func (u *User) PasswordHash() []byte {
	return u.passwordHash
}
